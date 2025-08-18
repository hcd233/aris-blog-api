package recommendation

import (
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/hcd233/aris-blog-api/internal/constant"
)

// Matrix 表示用户-物品评分矩阵
type Matrix struct {
	Data    map[int]map[int]float64 // [userID][itemID] -> rating
	Users   []int                   // 用户ID列表
	Items   []int                   // 物品ID列表
	UserMap map[int]int             // userID -> index
	ItemMap map[int]int             // itemID -> index
	mutex   sync.RWMutex
}

// NewMatrix 创建新的矩阵
func NewMatrix() *Matrix {
	return &Matrix{
		Data:    make(map[int]map[int]float64),
		UserMap: make(map[int]int),
		ItemMap: make(map[int]int),
	}
}

// Set 设置评分
func (m *Matrix) Set(userID, itemID int, rating float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.Data[userID] == nil {
		m.Data[userID] = make(map[int]float64)
	}
	m.Data[userID][itemID] = rating

	// 更新用户和物品索引
	if _, exists := m.UserMap[userID]; !exists {
		m.UserMap[userID] = len(m.Users)
		m.Users = append(m.Users, userID)
	}
	if _, exists := m.ItemMap[itemID]; !exists {
		m.ItemMap[itemID] = len(m.Items)
		m.Items = append(m.Items, itemID)
	}
}

// Get 获取评分
func (m *Matrix) Get(userID, itemID int) float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if userRatings, exists := m.Data[userID]; exists {
		if rating, exists := userRatings[itemID]; exists {
			return rating
		}
	}
	return 0.0
}

// GetUserRatings 获取用户的所有评分
func (m *Matrix) GetUserRatings(userID int) map[int]float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if userRatings, exists := m.Data[userID]; exists {
		result := make(map[int]float64)
		for itemID, rating := range userRatings {
			result[itemID] = rating
		}
		return result
	}
	return make(map[int]float64)
}

// MatrixFactorization 矩阵分解结构
type MatrixFactorization struct {
	UserFeatures [][]float64 // 用户特征矩阵 [users][factors]
	ItemFeatures [][]float64 // 物品特征矩阵 [items][factors]
	UserBias     []float64   // 用户偏置
	ItemBias     []float64   // 物品偏置
	GlobalMean   float64     // 全局均值
	Factors      int         // 因子数量
	mutex        sync.RWMutex
}

// NewMatrixFactorization 创建矩阵分解实例
func NewMatrixFactorization(userCount, itemCount, factors int) *MatrixFactorization {
	rand.Seed(time.Now().UnixNano())

	mf := &MatrixFactorization{
		UserFeatures: make([][]float64, userCount),
		ItemFeatures: make([][]float64, itemCount),
		UserBias:     make([]float64, userCount),
		ItemBias:     make([]float64, itemCount),
		Factors:      factors,
	}

	// 初始化特征矩阵
	for i := 0; i < userCount; i++ {
		mf.UserFeatures[i] = make([]float64, factors)
		for j := 0; j < factors; j++ {
			mf.UserFeatures[i][j] = rand.Float64()*0.1 - 0.05 // [-0.05, 0.05]
		}
	}

	for i := 0; i < itemCount; i++ {
		mf.ItemFeatures[i] = make([]float64, factors)
		for j := 0; j < factors; j++ {
			mf.ItemFeatures[i][j] = rand.Float64()*0.1 - 0.05 // [-0.05, 0.05]
		}
	}

	return mf
}

// Predict 预测评分
func (mf *MatrixFactorization) Predict(userIdx, itemIdx int) float64 {
	mf.mutex.RLock()
	defer mf.mutex.RUnlock()

	if userIdx >= len(mf.UserFeatures) || itemIdx >= len(mf.ItemFeatures) {
		return mf.GlobalMean
	}

	prediction := mf.GlobalMean + mf.UserBias[userIdx] + mf.ItemBias[itemIdx]

	// 计算特征向量点积
	for f := 0; f < mf.Factors; f++ {
		prediction += mf.UserFeatures[userIdx][f] * mf.ItemFeatures[itemIdx][f]
	}

	// 限制预测值范围
	if prediction < constant.CFMinRating {
		prediction = constant.CFMinRating
	}
	if prediction > constant.CFMaxRating {
		prediction = constant.CFMaxRating
	}

	return prediction
}

// Train 训练矩阵分解模型
func (mf *MatrixFactorization) Train(matrix *Matrix, learningRate, regularization float64, iterations int) error {
	if len(matrix.Users) == 0 || len(matrix.Items) == 0 {
		return nil
	}

	// 计算全局均值
	var totalRating float64
	var ratingCount int
	for _, userRatings := range matrix.Data {
		for _, rating := range userRatings {
			totalRating += rating
			ratingCount++
		}
	}
	if ratingCount > 0 {
		mf.GlobalMean = totalRating / float64(ratingCount)
	}

	// SGD训练
	for iter := 0; iter < iterations; iter++ {
		for userID, userRatings := range matrix.Data {
			userIdx := matrix.UserMap[userID]
			for itemID, actualRating := range userRatings {
				itemIdx := matrix.ItemMap[itemID]

				// 预测评分
				prediction := mf.Predict(userIdx, itemIdx)
				error := actualRating - prediction

				// 更新偏置
				userBiasOld := mf.UserBias[userIdx]
				itemBiasOld := mf.ItemBias[itemIdx]

				mf.UserBias[userIdx] += learningRate * (error - regularization*userBiasOld)
				mf.ItemBias[itemIdx] += learningRate * (error - regularization*itemBiasOld)

				// 更新特征向量
				for f := 0; f < mf.Factors; f++ {
					userFeatureOld := mf.UserFeatures[userIdx][f]
					itemFeatureOld := mf.ItemFeatures[itemIdx][f]

					mf.UserFeatures[userIdx][f] += learningRate * (error*itemFeatureOld - regularization*userFeatureOld)
					mf.ItemFeatures[itemIdx][f] += learningRate * (error*userFeatureOld - regularization*itemFeatureOld)
				}
			}
		}

		// 检查收敛
		if iter%10 == 0 {
			rmse := mf.calculateRMSE(matrix)
			if rmse < constant.CFConvergenceThreshold {
				break
			}
		}
	}

	return nil
}

// calculateRMSE 计算均方根误差
func (mf *MatrixFactorization) calculateRMSE(matrix *Matrix) float64 {
	var totalError float64
	var count int

	for userID, userRatings := range matrix.Data {
		userIdx := matrix.UserMap[userID]
		for itemID, actualRating := range userRatings {
			itemIdx := matrix.ItemMap[itemID]
			prediction := mf.Predict(userIdx, itemIdx)
			error := actualRating - prediction
			totalError += error * error
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return math.Sqrt(totalError / float64(count))
}

// RecommendationResult 推荐结果
type RecommendationResult struct {
	ItemID int     `json:"itemId"`
	Score  float64 `json:"score"`
	Reason string  `json:"reason"`
}

// CollaborativeFilter 协同过滤推荐器
type CollaborativeFilter struct {
	matrix *Matrix
	mf     *MatrixFactorization
	mutex  sync.RWMutex
}

// NewCollaborativeFilter 创建协同过滤推荐器
func NewCollaborativeFilter() *CollaborativeFilter {
	return &CollaborativeFilter{
		matrix: NewMatrix(),
	}
}

// AddRating 添加评分数据
func (cf *CollaborativeFilter) AddRating(userID, itemID int, rating float64) {
	cf.matrix.Set(userID, itemID, rating)
}

// Train 训练模型
func (cf *CollaborativeFilter) Train() error {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	userCount := len(cf.matrix.Users)
	itemCount := len(cf.matrix.Items)

	if userCount == 0 || itemCount == 0 {
		return nil
	}

	cf.mf = NewMatrixFactorization(userCount, itemCount, constant.CFDefaultFactors)

	return cf.mf.Train(
		cf.matrix,
		constant.CFDefaultLearningRate,
		constant.CFDefaultRegularization,
		constant.CFDefaultIterations,
	)
}

// Recommend 生成推荐
func (cf *CollaborativeFilter) Recommend(userID int, excludeItems []int, limit int) []RecommendationResult {
	cf.mutex.RLock()
	defer cf.mutex.RUnlock()

	if cf.mf == nil {
		return []RecommendationResult{}
	}

	userIdx, exists := cf.matrix.UserMap[userID]
	if !exists {
		return []RecommendationResult{}
	}

	// 获取用户已评分的物品
	userRatings := cf.matrix.GetUserRatings(userID)
	excludeSet := make(map[int]bool)
	for itemID := range userRatings {
		excludeSet[itemID] = true
	}
	for _, itemID := range excludeItems {
		excludeSet[itemID] = true
	}

	// 生成推荐
	var recommendations []RecommendationResult
	for _, itemID := range cf.matrix.Items {
		if excludeSet[itemID] {
			continue
		}

		itemIdx := cf.matrix.ItemMap[itemID]
		score := cf.mf.Predict(userIdx, itemIdx)

		recommendations = append(recommendations, RecommendationResult{
			ItemID: itemID,
			Score:  score,
			Reason: "协同过滤推荐",
		})
	}

	// 按分数排序
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// 限制返回数量
	if limit > 0 && len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations
}

// GetUserSimilarity 计算用户相似度
func (cf *CollaborativeFilter) GetUserSimilarity(userID1, userID2 int) float64 {
	cf.mutex.RLock()
	defer cf.mutex.RUnlock()

	if cf.mf == nil {
		return 0.0
	}

	userIdx1, exists1 := cf.matrix.UserMap[userID1]
	userIdx2, exists2 := cf.matrix.UserMap[userID2]

	if !exists1 || !exists2 {
		return 0.0
	}

	// 计算余弦相似度
	var dotProduct, norm1, norm2 float64
	for f := 0; f < cf.mf.Factors; f++ {
		feature1 := cf.mf.UserFeatures[userIdx1][f]
		feature2 := cf.mf.UserFeatures[userIdx2][f]
		dotProduct += feature1 * feature2
		norm1 += feature1 * feature1
		norm2 += feature2 * feature2
	}

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

// GetItemSimilarity 计算物品相似度
func (cf *CollaborativeFilter) GetItemSimilarity(itemID1, itemID2 int) float64 {
	cf.mutex.RLock()
	defer cf.mutex.RUnlock()

	if cf.mf == nil {
		return 0.0
	}

	itemIdx1, exists1 := cf.matrix.ItemMap[itemID1]
	itemIdx2, exists2 := cf.matrix.ItemMap[itemID2]

	if !exists1 || !exists2 {
		return 0.0
	}

	// 计算余弦相似度
	var dotProduct, norm1, norm2 float64
	for f := 0; f < cf.mf.Factors; f++ {
		feature1 := cf.mf.ItemFeatures[itemIdx1][f]
		feature2 := cf.mf.ItemFeatures[itemIdx2][f]
		dotProduct += feature1 * feature2
		norm1 += feature1 * feature1
		norm2 += feature2 * feature2
	}

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}