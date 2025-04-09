package etl

import (
	"fmt"
	"io/ioutil"
	"meta_commerce/app/etl/dest"
	"meta_commerce/app/etl/source"
	"meta_commerce/app/etl/transform"
	"meta_commerce/app/etl/types"
	"os"
	"strings"
	"time"

	"meta_commerce/app/registry"

	"gopkg.in/yaml.v3"
)

// PipelineConfig chứa cấu hình đầy đủ cho một pipeline
type PipelineConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Owner       string `yaml:"owner"`

	Source      SourceConfig      `yaml:"source"`
	Transform   []TransformConfig `yaml:"transform"`
	Destination DestConfig        `yaml:"destination"`
	Schedule    []ScheduleConfig  `yaml:"schedule"`
	Validation  ValidationConfig  `yaml:"validation"`
}

// SourceConfig cấu hình cho source
type SourceConfig struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

// TransformConfig cấu hình cho một bước transform
type TransformConfig struct {
	Name   string                 `yaml:"name"`
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

// DestConfig cấu hình cho destination
type DestConfig struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

// ScheduleConfig cấu hình lịch chạy
type ScheduleConfig struct {
	Name    string        `yaml:"name"`
	Cron    string        `yaml:"cron"`
	Enabled bool          `yaml:"enabled"`
	Timeout time.Duration `yaml:"timeout"`
}

// ValidationConfig cấu hình validation
type ValidationConfig struct {
	PreConditions  []ValidationCheck `yaml:"pre_conditions"`
	PostConditions []ValidationCheck `yaml:"post_conditions"`
}

// ValidationCheck cấu hình một validation check
type ValidationCheck struct {
	Check   string                 `yaml:"check"`
	Timeout string                 `yaml:"timeout"`
	Params  map[string]interface{} `yaml:"params"`
}

// PipelineLoader load và khởi tạo pipeline từ config
type PipelineLoader struct {
	configPath string
	registry   *registry.ETLRegistry
}

// NewPipelineLoader tạo một instance mới của PipelineLoader
func NewPipelineLoader(configPath string) *PipelineLoader {
	return &PipelineLoader{
		configPath: configPath,
		registry:   registry.GetETLRegistry(),
	}
}

// LoadPipeline load pipeline từ file config
func (l *PipelineLoader) LoadPipeline(name string) (*Pipeline, error) {
	// Đọc file config
	data, err := ioutil.ReadFile(l.configPath)
	if err != nil {
		return nil, fmt.Errorf("lỗi đọc file config: %v", err)
	}

	// Parse YAML
	var config struct {
		Pipelines map[string]PipelineConfig `yaml:"pipelines"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("lỗi parse YAML: %v", err)
	}

	// Lấy config của pipeline cần load
	pipelineConfig, exists := config.Pipelines[name]
	if !exists {
		return nil, fmt.Errorf("không tìm thấy pipeline %s", name)
	}

	// 1. Tạo source từ registry
	src, err := l.registry.CreateSource(pipelineConfig.Source.Type, pipelineConfig.Source.Config)
	if err != nil {
		return nil, fmt.Errorf("lỗi tạo source: %v", err)
	}

	// 2. Tạo transformer từ registry
	var transformer types.Transformer
	for _, config := range pipelineConfig.Transform {
		if config.Type == "field_mapper" {
			transformer, err = l.registry.CreateTransformer(config.Type, config.Config)
			if err != nil {
				return nil, fmt.Errorf("lỗi tạo transformer: %v", err)
			}
			break
		}
	}
	if transformer == nil {
		return nil, fmt.Errorf("không tìm thấy field mapper trong config")
	}

	// 3. Tạo destination từ registry
	dst, err := l.registry.CreateDestination(pipelineConfig.Destination.Type, pipelineConfig.Destination.Config)
	if err != nil {
		return nil, fmt.Errorf("lỗi tạo destination: %v", err)
	}

	// 4. Tạo pipeline
	pipeline, err := NewPipeline(src, transformer, dst)
	if err != nil {
		return nil, fmt.Errorf("lỗi tạo pipeline: %v", err)
	}

	return pipeline, nil
}

// createSource tạo source component từ config
func (l *PipelineLoader) createSource(config SourceConfig) (types.DataSource, error) {
	switch config.Type {
	case "rest_api":
		// Parse config
		url := config.Config["url"].(string)
		timeout := config.Config["timeout"].(string)
		timeoutDuration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, fmt.Errorf("lỗi parse timeout: %v", err)
		}

		// Xử lý environment variables trong headers
		headers := make(map[string]string)
		if headersConfig, ok := config.Config["headers"].(map[string]interface{}); ok {
			for k, v := range headersConfig {
				value := v.(string)
				if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
					envVar := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
					value = os.Getenv(envVar)
				}
				headers[k] = value
			}
		}

		// Tạo source config
		sourceConfig := source.RestConfig{
			URL:     url,
			Method:  config.Config["method"].(string),
			Headers: headers,
			Timeout: timeoutDuration,
		}

		return source.NewRestApiSource(sourceConfig)

	default:
		return nil, fmt.Errorf("không hỗ trợ source type: %s", config.Type)
	}
}

// createTransformer tạo transformer từ config
func (l *PipelineLoader) createTransformer(configs []TransformConfig) (types.Transformer, error) {
	// Hiện tại chỉ hỗ trợ field mapper
	for _, config := range configs {
		if config.Type == "field_mapper" {
			// Parse mappings
			var mappings []transform.FieldMapping
			mappingsConfig := config.Config["mappings"].([]interface{})
			for _, m := range mappingsConfig {
				mapping := m.(map[string]interface{})
				mappings = append(mappings, transform.FieldMapping{
					Source: mapping["source"].(string),
					Target: mapping["target"].(string),
					Type:   mapping["type"].(string),
				})
			}

			// Tạo transformer config
			transformerConfig := transform.MapperConfig{
				Mappings: mappings,
			}

			return transform.NewFieldMapper(transformerConfig)
		}
	}

	return nil, fmt.Errorf("không tìm thấy field mapper trong config")
}

// createDestination tạo destination từ config
func (l *PipelineLoader) createDestination(config DestConfig) (types.Destination, error) {
	switch config.Type {
	case "internal_api":
		// Parse config
		url := config.Config["url"].(string)
		timeout := config.Config["timeout"].(string)
		timeoutDuration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, fmt.Errorf("lỗi parse timeout: %v", err)
		}

		// Xử lý environment variables trong headers
		headers := make(map[string]string)
		if headersConfig, ok := config.Config["headers"].(map[string]interface{}); ok {
			for k, v := range headersConfig {
				value := v.(string)
				if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
					envVar := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
					value = os.Getenv(envVar)
				}
				headers[k] = value
			}
		}

		// Tạo destination config
		destConfig := dest.InternalAPIConfig{
			URL:     url,
			Method:  config.Config["method"].(string),
			Headers: headers,
			Timeout: timeoutDuration,
		}

		return dest.NewInternalAPIDestination(destConfig)

	default:
		return nil, fmt.Errorf("không hỗ trợ destination type: %s", config.Type)
	}
}
