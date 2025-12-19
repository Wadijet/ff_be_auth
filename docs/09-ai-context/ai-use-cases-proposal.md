# Đề Xuất Ứng Dụng AI Cho Tập Dữ Liệu Folkform

## 📋 Mục Lục

1. [Tổng Quan Dữ Liệu](#tổng-quan-dữ-liệu)
2. [Use Cases Ưu Tiên Cao](#use-cases-ưu-tiên-cao)
3. [Use Cases Ưu Tiên Trung Bình](#use-cases-ưu-tiên-trung-bình)
4. [Use Cases Dài Hạn](#use-cases-dài-hạn)
5. [Roadmap Triển Khai](#roadmap-triển-khai)
6. [Technical Requirements](#technical-requirements)

---

## Tổng Quan Dữ Liệu

### Dữ Liệu Hiện Có

| Collection | Số Lượng | Đặc Điểm | AI Potential |
|-----------|----------|----------|--------------|
| **customers** | 33,110 | Multi-source (Pancake + POS), có phone, name, gender | ⭐⭐⭐⭐⭐ |
| **pc_pos_orders** | 2,633 | Có customerId, pageId, postId, status | ⭐⭐⭐⭐⭐ |
| **fb_message_items** | 834,756 | Text messages, attachments, timestamps | ⭐⭐⭐⭐⭐ |
| **fb_conversations** | 26,832 | Metadata, tags, assignees, ad_ids | ⭐⭐⭐⭐ |
| **fb_posts** | 5,249 | Posts với reactions, comments | ⭐⭐⭐ |
| **pc_pos_products** | 401 | Products với attributes, variations | ⭐⭐⭐ |
| **pc_pos_variations** | 2,820 | Variations với price, quantity, images | ⭐⭐⭐ |

### Đặc Điểm Dữ Liệu

✅ **Strengths:**
- Lượng dữ liệu lớn (834K messages)
- Dữ liệu đa dạng (text, structured, timestamps)
- Có mối quan hệ rõ ràng (customer → conversation → order)
- Có timestamps đầy đủ cho time-series analysis

⚠️ **Limitations:**
- Order items chưa được extract (cần fix trước)
- Shipping address chưa được extract
- Customer POS data chưa sync
- Một số fields quan trọng còn thiếu

---

## Use Cases Ưu Tiên Cao

### 1. 🤖 AI Customer Service Assistant

**Mô tả**: Chatbot tự động trả lời tin nhắn khách hàng trên Facebook, giảm workload cho CS team.

**Input Data:**
- `fb_message_items.messageData.message` (834K messages)
- `fb_message_items.messageData.attachments` (images, files)
- `fb_conversations.panCakeData.snippet` (context)
- `customers` (thông tin khách hàng)
- `pc_pos_products` (thông tin sản phẩm)

**AI Model:**
- **LLM**: GPT-4, Claude, hoặc Vietnamese LLM (VinAI, PhoBERT-based)
- **RAG**: Vector database (Pinecone, Weaviate) chứa:
  - Product catalog
  - FAQ
  - Previous conversations
  - Order history

**Output:**
- Câu trả lời tự động phù hợp
- Confidence score
- Suggested actions (tạo đơn, chuyển human, etc.)

**Value Proposition:**
- ⏱️ **Giảm 60-80% response time** (từ vài giờ → vài phút)
- 💰 **Tiết kiệm 40-60% chi phí CS** (tự động hóa 70% queries đơn giản)
- 📈 **Tăng customer satisfaction** (phản hồi nhanh, 24/7)

**Implementation:**
```python
# Pseudo-code
def ai_customer_service(message, customer_id, conversation_id):
    # 1. Retrieve context
    context = get_conversation_context(conversation_id)
    customer = get_customer(customer_id)
    products = search_products(message)  # Semantic search
    
    # 2. Classify intent
    intent = classify_intent(message)  # Hỏi giá, khiếu nại, đặt hàng, etc.
    
    # 3. Generate response
    if intent == "product_inquiry":
        response = generate_product_response(message, products, customer)
    elif intent == "order_status":
        response = get_order_status(customer)
    elif intent == "complaint":
        response = escalate_to_human(message, customer)
    else:
        response = llm.generate(message, context, customer, products)
    
    return response, intent, confidence_score
```

**Metrics:**
- Response accuracy: >85%
- Customer satisfaction: >4.0/5.0
- Escalation rate: <20%
- Cost per conversation: -60%

**Priority**: 🔴 **HIGH** - Impact lớn, ROI cao, dữ liệu đã sẵn sàng

---

### 2. 📊 Sentiment Analysis & Customer Satisfaction

**Mô tả**: Phân tích sentiment của messages để phát hiện khách hàng không hài lòng sớm và can thiệp.

**Input Data:**
- `fb_message_items.messageData.message` (834K messages)
- `fb_message_items.messageData.from` (sender info)
- `fb_conversations.panCakeData.tag_histories` (tags)
- `pc_pos_orders.status` (order status - nếu có)

**AI Model:**
- **Sentiment Analysis**: 
  - Vietnamese BERT models (PhoBERT, vBERT)
  - Fine-tuned trên dataset customer service
- **Emotion Detection**: Multi-label classification (angry, happy, neutral, frustrated)

**Output:**
- Sentiment score: -1 (negative) → +1 (positive)
- Emotion labels: [angry, frustrated, happy, neutral, etc.]
- Risk score: 0-100 (khả năng churn/complaint)
- Alert khi sentiment < threshold

**Value Proposition:**
- 🚨 **Phát hiện sớm 80% complaints** trước khi escalate
- 📈 **Tăng NPS 15-20%** (can thiệp sớm)
- 💰 **Giảm 30% refund/return** (giải quyết vấn đề sớm)

**Implementation:**
```python
def analyze_sentiment(conversation_id):
    messages = get_messages(conversation_id)
    
    # Analyze each message
    sentiments = []
    for msg in messages:
        sentiment = sentiment_model.predict(msg['message'])
        emotion = emotion_model.predict(msg['message'])
        sentiments.append({
            'sentiment': sentiment,
            'emotion': emotion,
            'timestamp': msg['insertedAt']
        })
    
    # Aggregate conversation sentiment
    avg_sentiment = mean([s['sentiment'] for s in sentiments])
    risk_score = calculate_risk(sentiments, customer_history)
    
    # Alert if negative
    if avg_sentiment < -0.3 or risk_score > 70:
        alert_cs_team(conversation_id, risk_score, sentiments)
    
    return {
        'conversation_sentiment': avg_sentiment,
        'risk_score': risk_score,
        'emotions': aggregate_emotions(sentiments),
        'trend': sentiment_trend(sentiments)  # improving/worsening
    }
```

**Metrics:**
- Sentiment accuracy: >90%
- Early detection rate: >80%
- False positive rate: <10%

**Priority**: 🔴 **HIGH** - Dữ liệu sẵn sàng, impact cao

---

### 3. 🎯 Lead Scoring & Conversion Prediction

**Mô tả**: Dự đoán khả năng khách hàng chuyển đổi từ conversation → order.

**Input Data:**
- `fb_conversations` (26K conversations)
- `fb_message_items` (messages trong conversation)
- `customers` (customer profile)
- `pc_pos_orders` (historical orders)
- `fb_posts` (nếu conversation từ post)

**AI Model:**
- **Classification**: XGBoost, LightGBM, hoặc Neural Network
- **Features**:
  - Conversation features: message count, response time, sentiment, intent
  - Customer features: total orders, total spent, last order date
  - Engagement features: post engagement, ad clicks
  - Temporal features: time of day, day of week

**Output:**
- Conversion probability: 0-100%
- Lead score: 0-100
- Time to convert prediction: X days
- Recommended actions: [follow_up, send_promotion, assign_sales]

**Value Proposition:**
- 📈 **Tăng conversion rate 25-40%** (focus vào high-quality leads)
- ⏱️ **Giảm sales cycle 30%** (prioritize hot leads)
- 💰 **Tăng revenue 20-30%** (better lead qualification)

**Implementation:**
```python
def predict_conversion(conversation_id):
    conversation = get_conversation(conversation_id)
    customer = get_customer(conversation['customerId'])
    messages = get_messages(conversation_id)
    
    # Extract features
    features = {
        # Conversation features
        'message_count': len(messages),
        'avg_response_time': calculate_avg_response_time(messages),
        'sentiment': analyze_sentiment(messages),
        'intent': classify_intent(messages),
        'has_product_inquiry': check_product_inquiry(messages),
        'has_price_inquiry': check_price_inquiry(messages),
        
        # Customer features
        'customer_total_orders': customer.get('totalOrder', 0),
        'customer_total_spent': customer.get('totalSpent', 0),
        'customer_last_order_days_ago': days_since_last_order(customer),
        'customer_is_returning': customer.get('totalOrder', 0) > 0,
        
        # Engagement features
        'conversation_duration_hours': calculate_duration(conversation),
        'messages_per_hour': len(messages) / conversation_duration_hours,
        'has_attachment': any(msg.get('attachments') for msg in messages),
        
        # Temporal features
        'hour_of_day': extract_hour(conversation['createdAt']),
        'day_of_week': extract_day_of_week(conversation['createdAt']),
    }
    
    # Predict
    conversion_prob = conversion_model.predict_proba(features)[1]
    lead_score = calculate_lead_score(features, conversion_prob)
    time_to_convert = time_to_convert_model.predict(features)
    
    # Recommend actions
    actions = recommend_actions(conversation, customer, conversion_prob, lead_score)
    
    return {
        'conversation_id': conversation_id,
        'conversion_probability': conversion_prob,
        'lead_score': lead_score,
        'time_to_convert_days': time_to_convert,
        'recommended_actions': actions,
        'key_factors': explain_prediction(features)  # Why this score?
    }
```

**Metrics:**
- Prediction accuracy: >80%
- Precision@Top20%: >60% (60% of top 20% actually convert)
- ROI: 3-5x (revenue increase / AI cost)

**Priority**: 🔴 **HIGH** - Direct impact on revenue

---

### 4. 📦 Product Recommendation Engine

**Mô tả**: Gợi ý sản phẩm phù hợp cho khách hàng dựa trên lịch sử mua hàng, conversations, và preferences.

**Input Data:**
- `pc_pos_orders.orderItems` (sau khi extract) - purchase history
- `fb_message_items.messageData.message` - product inquiries
- `customers` - customer profile, preferences
- `pc_pos_products` - product catalog với attributes
- `pc_pos_variations` - variations với images, prices

**AI Model:**
- **Collaborative Filtering**: Matrix factorization (users × products)
- **Content-Based Filtering**: Product attributes matching
- **Hybrid**: Combine both approaches
- **Deep Learning**: Neural Collaborative Filtering (NCF)

**Output:**
- Top N recommended products với scores
- Explanation: "Vì bạn đã mua X, bạn có thể thích Y"
- Personalized product bundles

**Value Proposition:**
- 📈 **Tăng cross-sell 30-50%** (gợi ý sản phẩm liên quan)
- 💰 **Tăng AOV 15-25%** (upsell, bundles)
- 🎯 **Tăng conversion 20%** (relevant recommendations)

**Implementation:**
```python
def recommend_products(customer_id, context=None):
    customer = get_customer(customer_id)
    order_history = get_order_history(customer_id)
    conversations = get_recent_conversations(customer_id)
    
    # Extract preferences from conversations
    preferences = extract_preferences(conversations)  # colors, styles, price range
    
    # Collaborative filtering
    cf_recommendations = collaborative_filtering.recommend(
        customer_id, 
        order_history,
        n_recommendations=10
    )
    
    # Content-based filtering
    cb_recommendations = content_based.recommend(
        customer_preferences=preferences,
        order_history=order_history,
        product_catalog=products,
        n_recommendations=10
    )
    
    # Hybrid approach
    final_recommendations = hybrid_recommend(
        cf_recommendations,
        cb_recommendations,
        weights=[0.6, 0.4]  # CF 60%, CB 40%
    )
    
    # Add explanations
    for rec in final_recommendations:
        rec['explanation'] = generate_explanation(rec, customer, order_history)
        rec['bundle_suggestions'] = find_bundles(rec['product_id'])
    
    return final_recommendations
```

**Metrics:**
- Recommendation accuracy (precision@10): >40%
- Click-through rate: >15%
- Conversion rate: >5%
- Revenue lift: +20-30%

**Priority**: 🟡 **MEDIUM** - Cần extract orderItems trước

---

### 5. 🔮 Churn Prediction & Retention

**Mô tả**: Dự đoán khách hàng có nguy cơ rời bỏ và đề xuất actions để giữ chân.

**Input Data:**
- `customers.posLastOrderAt` (sau khi sync POS)
- `customers.totalOrder`, `totalSpent`
- `pc_pos_orders.insertedAt` (order frequency)
- `fb_conversations` (engagement level)
- `fb_message_items` (sentiment, interaction)

**AI Model:**
- **Classification**: XGBoost, Random Forest
- **Survival Analysis**: Cox Proportional Hazards (time to churn)
- **Features**:
  - Recency: Days since last order
  - Frequency: Order count
  - Monetary: Total spent
  - Engagement: Conversation count, message count, sentiment
  - Product diversity: Number of unique products bought

**Output:**
- Churn probability: 0-100%
- Churn risk level: Low/Medium/High
- Days until predicted churn: X days
- Recommended retention actions: [discount, new_product, re_engagement_campaign]

**Value Proposition:**
- 💰 **Giảm churn rate 25-40%** (can thiệp sớm)
- 📈 **Tăng LTV 20-30%** (giữ chân khách hàng)
- 🎯 **ROI retention campaigns: 5-10x**

**Implementation:**
```python
def predict_churn(customer_id):
    customer = get_customer(customer_id)
    orders = get_orders(customer_id)
    conversations = get_conversations(customer_id)
    
    # Calculate features
    features = {
        'recency_days': days_since_last_order(customer, orders),
        'frequency': len(orders),
        'monetary': customer.get('totalSpent', 0),
        'avg_order_value': customer.get('totalSpent', 0) / max(len(orders), 1),
        'order_frequency_days': calculate_order_frequency(orders),
        'conversation_count': len(conversations),
        'last_conversation_days_ago': days_since_last_conversation(conversations),
        'avg_sentiment': calculate_avg_sentiment(conversations),
        'product_diversity': count_unique_products(orders),
        'return_rate': calculate_return_rate(orders),
    }
    
    # Predict
    churn_prob = churn_model.predict_proba(features)[1]
    days_to_churn = survival_model.predict(features)
    risk_level = classify_risk(churn_prob, days_to_churn)
    
    # Recommend actions
    actions = recommend_retention_actions(customer, churn_prob, risk_level)
    
    return {
        'customer_id': customer_id,
        'churn_probability': churn_prob,
        'risk_level': risk_level,
        'days_to_churn': days_to_churn,
        'recommended_actions': actions,
        'key_factors': explain_churn_risk(features)
    }
```

**Metrics:**
- Prediction accuracy: >75%
- Precision@HighRisk: >60%
- Retention rate improvement: +25-40%

**Priority**: 🟡 **MEDIUM** - Cần sync POS customer data trước

---

## Use Cases Ưu Tiên Trung Bình

### 6. 📝 Intent Classification & Auto-Routing

**Mô tả**: Tự động phân loại intent của messages và route đến đúng bộ phận/handler.

**Input Data:**
- `fb_message_items.messageData.message`
- `fb_conversations.panCakeData.type` (INBOX, COMMENT, LIVESTREAM)

**AI Model:**
- **Text Classification**: BERT-based (PhoBERT fine-tuned)
- **Intent Labels**: 
  - Product inquiry
  - Price inquiry
  - Order status
  - Complaint
  - Return/Refund
  - General question
  - Spam

**Output:**
- Intent label với confidence
- Suggested handler: [sales, cs, logistics, etc.]
- Auto-response template (nếu có)

**Value Proposition:**
- ⏱️ **Giảm 50% routing time**
- 📈 **Tăng 30% first-response accuracy**
- 💰 **Giảm 20% CS workload**

**Priority**: 🟡 **MEDIUM**

---

### 7. 💬 Conversation Summarization

**Mô tả**: Tự động tóm tắt conversations dài để CS team nắm nhanh context.

**Input Data:**
- `fb_message_items` (tất cả messages trong conversation)
- `fb_conversations.panCakeData`

**AI Model:**
- **Summarization**: BART, T5 (Vietnamese fine-tuned)
- **Extractive + Abstractive**: Combine both approaches

**Output:**
- Conversation summary (2-3 sentences)
- Key points: [main_issue, customer_request, resolution_status]
- Action items: [follow_up_needed, order_to_create, etc.]

**Value Proposition:**
- ⏱️ **Giảm 70% time để hiểu context**
- 📈 **Tăng 40% CS efficiency**

**Priority**: 🟡 **MEDIUM**

---

### 8. 📊 Sales Forecasting

**Mô tả**: Dự báo doanh thu, số đơn hàng trong tương lai.

**Input Data:**
- `pc_pos_orders.insertedAt` (time series)
- `pc_pos_orders.total_price` (sau khi extract)
- `fb_conversations` (lead pipeline)
- `fb_posts` (marketing activities)

**AI Model:**
- **Time Series**: Prophet, ARIMA, LSTM, Transformer-based (Temporal Fusion Transformer)
- **Features**: Historical sales, seasonality, trends, external factors

**Output:**
- Daily/Weekly/Monthly revenue forecast
- Order count forecast
- Confidence intervals
- Anomaly detection

**Value Proposition:**
- 📈 **Cải thiện inventory planning**
- 💰 **Tối ưu marketing budget**
- 🎯 **Dự báo chính xác ±10%**

**Priority**: 🟡 **MEDIUM** - Cần extract total_price trước

---

### 9. 🖼️ Image Analysis for Product Recommendations

**Mô tả**: Phân tích images trong messages để hiểu customer preferences và recommend products.

**Input Data:**
- `fb_message_items.messageData.attachments` (images)
- `pc_pos_products`, `pc_pos_variations` (product images)

**AI Model:**
- **Image Classification**: ResNet, EfficientNet
- **Similarity Search**: CLIP (text-image matching)
- **Style/Color Detection**: Computer vision models

**Output:**
- Detected style/color preferences
- Similar products (visual similarity)
- Product recommendations based on images

**Value Proposition:**
- 📈 **Tăng 25% conversion** (visual matching)
- 🎯 **Better understanding customer taste**

**Priority**: 🟢 **LOW** - Nice to have

---

## Use Cases Dài Hạn

### 10. 🤝 Customer Matching (Pancake ↔ POS)

**Mô tả**: Tự động match customers giữa Pancake (Facebook) và POS để có unified view.

**Input Data:**
- `customers` từ Pancake (phone, name, psid)
- `pc_pos_orders.customer` (phone, name, email)

**AI Model:**
- **Entity Resolution**: Fuzzy matching, record linkage
- **ML-based Matching**: Siamese networks, embeddings

**Output:**
- Matched customer pairs với confidence score
- Unified customer profile

**Value Proposition:**
- 📊 **Unified customer view**
- 📈 **Better analytics và personalization**

**Priority**: 🟢 **LOW** - Cần sync POS customers trước

---

### 11. 📱 Dynamic Pricing Optimization

**Mô tả**: Tối ưu giá sản phẩm dựa trên demand, inventory, customer segments.

**Input Data:**
- `pc_pos_orders` (sales data)
- `pc_pos_variations.quantity` (inventory)
- `customers` (segments)
- `fb_conversations` (demand signals)

**AI Model:**
- **Reinforcement Learning**: Multi-armed bandit, Q-learning
- **Optimization**: Price elasticity models

**Output:**
- Optimal price recommendations
- Price change impact prediction

**Value Proposition:**
- 💰 **Tăng revenue 10-20%**
- 📈 **Tối ưu inventory turnover**

**Priority**: 🟢 **LOW** - Advanced use case

---

## Roadmap Triển Khai

### Phase 1: Quick Wins (1-2 tháng)

**Mục tiêu**: Implement các use cases có ROI cao, dữ liệu sẵn sàng

1. ✅ **Sentiment Analysis** (2 tuần)
   - Fine-tune Vietnamese BERT model
   - Build real-time pipeline
   - Dashboard cho CS team

2. ✅ **Intent Classification** (2 tuần)
   - Train classification model
   - Auto-routing system
   - Integration với CS workflow

3. ✅ **AI Customer Service (MVP)** (4 tuần)
   - RAG setup với product catalog
   - Basic chatbot
   - Human handoff logic

**Expected ROI**: 2-3x trong 3 tháng đầu

---

### Phase 2: Core Features (2-4 tháng)

**Mục tiêu**: Build các features core cho business growth

4. ✅ **Lead Scoring** (3 tuần)
   - Feature engineering
   - Model training
   - Integration với sales workflow

5. ✅ **Product Recommendation** (4 tuần)
   - Collaborative + Content-based
   - API integration
   - A/B testing framework

6. ✅ **Churn Prediction** (3 tuần)
   - Model training
   - Retention campaign automation

**Expected ROI**: 4-5x trong 6 tháng

---

### Phase 3: Advanced Features (4-6 tháng)

**Mục tiêu**: Advanced AI features cho competitive advantage

7. ✅ **Sales Forecasting** (4 tuần)
8. ✅ **Conversation Summarization** (3 tuần)
9. ✅ **Image Analysis** (6 tuần)
10. ✅ **Dynamic Pricing** (8 tuần)

**Expected ROI**: 5-10x trong 12 tháng

---

## Technical Requirements

### Infrastructure

1. **ML Platform**:
   - Model training: AWS SageMaker, Google Vertex AI, hoặc self-hosted (MLflow)
   - Model serving: FastAPI, TensorFlow Serving, TorchServe
   - Vector database: Pinecone, Weaviate, hoặc Qdrant

2. **Data Pipeline**:
   - ETL: Apache Airflow, Prefect
   - Feature store: Feast, Tecton
   - Real-time: Kafka, Redis

3. **Monitoring**:
   - Model monitoring: Evidently AI, Fiddler
   - Performance tracking: MLflow, Weights & Biases

### Models & Libraries

1. **NLP**:
   - Vietnamese models: `vinai/phobert-base`, `FPTAI/vibert`
   - LLM: OpenAI GPT-4, Anthropic Claude, hoặc Vietnamese LLM
   - Libraries: `transformers`, `sentence-transformers`

2. **ML**:
   - `scikit-learn`, `xgboost`, `lightgbm`
   - `pytorch`, `tensorflow`

3. **Time Series**:
   - `prophet`, `statsmodels`, `pytorch-forecasting`

### Data Requirements

**Cần fix trước khi implement:**
1. ✅ Extract `orderItems` từ orders
2. ✅ Extract `shippingAddress` từ orders
3. ✅ Extract `total_price` từ orders
4. ✅ Sync POS customers
5. ✅ Populate `sources` field

---

## Kết Luận

Với **834K messages**, **33K customers**, và **2.6K orders**, hệ thống có đủ dữ liệu để implement các AI use cases có impact cao:

### Top 3 Use Cases Nên Bắt Đầu:

1. **🤖 AI Customer Service** - ROI cao nhất, giảm CS cost 40-60%
2. **📊 Sentiment Analysis** - Dữ liệu sẵn sàng, phát hiện complaints sớm
3. **🎯 Lead Scoring** - Direct impact on revenue, tăng conversion 25-40%

### Expected Overall Impact:

- 💰 **Revenue**: +20-30% trong 6 tháng
- ⏱️ **Efficiency**: +40-60% CS efficiency
- 📈 **Customer Satisfaction**: +15-20% NPS
- 💵 **Cost Savings**: -30-40% CS cost

### Next Steps:

1. **Fix data gaps** (orderItems, shippingAddress, etc.)
2. **Pilot Phase 1** (Sentiment Analysis + Intent Classification)
3. **Measure ROI** và scale up
4. **Iterate** dựa trên feedback

---

## Tài Liệu Tham Khảo

- [Data Structure Analysis](./data-structure-analysis.md)
- [Data Architecture Overview](./data-architecture-overview.md)
- [Pancake API Context](./pancake-api-context.md)
- [Pancake POS API Context](./pancake-pos-api-context.md)

