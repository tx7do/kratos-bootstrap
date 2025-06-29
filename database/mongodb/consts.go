package mongodb

const (
	// 比较操作符

	OperatorEq  = "$eq"  // 等于
	OperatorNe  = "$ne"  // 不等于
	OperatorGt  = "$gt"  // 大于
	OperatorGte = "$gte" // 大于等于
	OperatorLt  = "$lt"  // 小于
	OperatorLte = "$lte" // 小于等于

	// 逻辑操作符

	OperatorAnd = "$and" // 与
	OperatorOr  = "$or"  // 或
	OperatorNot = "$not" // 非
	OperatorNor = "$nor" // 非或

	// 元素操作符

	OperatorExists = "$exists" // 是否存在
	OperatorType   = "$type"   // 类型

	// 评估操作符

	OperatorExpr       = "$expr"       // 表达式
	OperatorJsonSchema = "$jsonSchema" // JSON Schema 验证
	OperatorMod        = "$mod"        // 取模
	OperatorRegex      = "$regex"      // 正则表达式
	OperatorText       = "$text"       // 文本搜索
	OperatorWhere      = "$where"      // JavaScript 表达式
	OperatorSearch     = "$search"     // 文本搜索

	// 数组操作符

	OperatorAll       = "$all"       // 匹配所有
	OperatorElemMatch = "$elemMatch" // 匹配数组中的元素
	OperatorSize      = "$size"      // 数组大小

	// 集合操作符

	OperatorIn  = "$in"  // 包含
	OperatorNin = "$nin" // 不包含

	// 更新操作符

	OperatorSet         = "$set"         // 设置字段值
	OperatorUnset       = "$unset"       // 删除字段
	OperatorInc         = "$inc"         // 增加值
	OperatorMul         = "$mul"         // 乘法
	OperatorRename      = "$rename"      // 重命名字段
	OperatorCurrentDate = "$currentDate" // 设置当前日期
	OperatorAddToSet    = "$addToSet"    // 添加到集合
	OperatorPop         = "$pop"         // 删除数组中的元素
	OperatorPull        = "$pull"        // 删除匹配的数组元素
	OperatorPush        = "$push"        // 添加数组元素
	OperatorEach        = "$each"        // 批量添加数组元素
	OperatorSlice       = "$slice"       // 截取数组
	OperatorSort        = "$sort"        // 排序数组
	OperatorPosition    = "$position"    // 指定数组位置

	// 聚合操作符

	OperatorGroup       = "$group"       // 分组
	OperatorMatch       = "$match"       // 匹配
	OperatorProject     = "$project"     // 投影
	OperatorSortAgg     = "$sort"        // 排序
	OperatorLimit       = "$limit"       // 限制
	OperatorSkip        = "$skip"        // 跳过
	OperatorUnwind      = "$unwind"      // 拆分数组
	OperatorLookup      = "$lookup"      // 关联查询
	OperatorAddFields   = "$addFields"   // 添加字段
	OperatorReplaceRoot = "$replaceRoot" // 替换根字段
	OperatorCount       = "$count"       // 计数
	OperatorFacet       = "$facet"       // 多面查询
	OperatorBucket      = "$bucket"      // 分桶
	OperatorBucketAuto  = "$bucketAuto"  // 自动分桶
	OperatorIndexStats  = "$indexStats"  // 索引统计
	OperatorOut         = "$out"         // 输出
	OperatorMerge       = "$merge"       // 合并
	OperatorSum         = "$sum"         // 求和
	OperatorAvg         = "$avg"         // 平均值
	OperatorMin         = "$min"         // 最小值
	OperatorMax         = "$max"         // 最大值
	OperatorFirst       = "$first"       // 第一个值
	OperatorLast        = "$last"        // 最后一个值
	OperatorStdDevPop   = "$stdDevPop"   // 总体标准差
	OperatorStdDevSamp  = "$stdDevSamp"  // 样本标准差

	// 类型转换操作符

	OperatorToLong    = "$toLong"    // 转换为 long 类型
	OperatorToDouble  = "$toDouble"  // 转换为 double 类型
	OperatorToDecimal = "$toDecimal" // 转换为 decimal 类型
	OperatorToString  = "$toString"  // 转换为 string 类型
	OperatorToDate    = "$toDate"    // 转换为 date 类型
	OperatorToInt     = "$toInt"     // 转换为 int 类型

	// 地理空间操作符

	OperatorNear          = "$near"          // 查询距离某点最近的文档
	OperatorNearSphere    = "$nearSphere"    // 查询距离某点最近的文档（球面距离）
	OperatorGeoWithin     = "$geoWithin"     // 地理范围查询
	OperatorGeoIntersects = "$geoIntersects" // 地理相交查询

	OperatorGeometry    = "$geometry"    // 几何图形
	OperatorMaxDistance = "$maxDistance" // 最大距离
	OperatorMinDistance = "$minDistance" // 最小距离
)
