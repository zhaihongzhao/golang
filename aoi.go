package aoi

// 位置
type Position struct {
	x, y float32
}

// 实体
type Entity interface {
	ID() int
	GetPos() Position
	SetPos(Position)
	OnEnterAOI(other Entity)
	OnLeaveAOI(other Entity)
}

// 网格
type Grid struct {
	entitys  map[int]Entity // 网格中的实体
	watchers map[int]Entity // 观察者实体
}

// 初始化
func (g *Grid) init() {
	g.entitys = make(map[int]Entity)
	g.watchers = make(map[int]Entity)
}

// 添加实体
func (g *Grid) addEntity(entity Entity) {
	g.entitys[entity.ID()] = entity
}

// 移除实体
func (g *Grid) removeEntity(entity Entity) {
	delete(g.entitys, entity.ID())
}

// 添加观察者
func (g *Grid) addWatcher(entity Entity) {
	if _, ok := g.watchers[entity.ID()]; ok {
		return
	}
	for _, other := range g.entitys {
		if entity.ID() != other.ID() {
			entity.OnEnterAOI(other)
		}
	}
	g.watchers[entity.ID()] = entity
}

// 移除观察者
func (g *Grid) removeWatcher(entity Entity) {
	if _, ok := g.watchers[entity.ID()]; !ok {
		return
	}
	for _, other := range g.entitys {
		if entity.ID() != other.ID() {
			entity.OnLeaveAOI(other)
		}
	}
	delete(g.watchers, entity.ID())
}

// 管理
type AOIManager struct {
	minX, maxX, minY, maxY float32  // 地图范围
	gsize                  float32  // 网格大小
	grids                  [][]Grid // 网格
	xNum, yNum             int      // 网格数量
}

// 创建管理
func NewAOIManager(minX, maxX, minY, maxY float32, gsize float32) *AOIManager {
	xNum := int((maxX-minX)/gsize) + 1
	yNum := int((maxY-minY)/gsize) + 1

	mgr := &AOIManager{
		minX:  minX,
		maxX:  maxX,
		minY:  minY,
		maxY:  maxY,
		gsize: gsize,
		xNum:  xNum,
		yNum:  yNum,
	}

	mgr.grids = make([][]Grid, xNum)
	for i := 0; i < xNum; i++ {
		mgr.grids[i] = make([]Grid, yNum)
		for j := 0; j < yNum; j++ {
			mgr.grids[i][j].init()
		}
	}

	return mgr
}

// 进入地图
func (m *AOIManager) Enter(entity Entity, pos Position) {
	// 添加实体
	entity.SetPos(pos)
	grid := m.posToGrid(pos)
	grid.addEntity(entity)
	// 添加观察者
	m.visitWatchGrids(pos, func(g *Grid) {
		g.addWatcher(entity)
	})
}

// 离开地图
func (m *AOIManager) Leave(entity Entity) {
	// 移除实体
	grid := m.posToGrid(entity.GetPos())
	grid.removeEntity(entity)
	// 移除观察者
	m.visitWatchGrids(entity.GetPos(), func(g *Grid) {
		g.removeWatcher(entity)
	})
}

// 移动
func (m *AOIManager) Move(entity Entity, toPos Position) {
	// 更新位置
	fromPos := entity.GetPos()
	entity.SetPos(toPos)
	fromGrid := m.posToGrid(fromPos)
	toGrid := m.posToGrid(toPos)
	if fromGrid == toGrid {
		return
	}

	// 跨越网格
	fromGrid.removeEntity(entity)
	toGrid.addEntity(entity)

	// 更新观察者
	fxmin, fxmax, fymin, fymax := m.getWatchGrids(fromPos)
	txmin, txmax, tymin, tymax := m.getWatchGrids(toPos)

	for x := fxmin; x <= fxmax; x++ {
		for y := fymin; y <= fymax; y++ {
			if x >= txmin && x <= txmax && y >= tymin && y <= tymax {
				continue
			}
			grid := &m.grids[x][y]
			grid.removeWatcher(entity)
		}
	}
	for x := txmin; x <= txmax; x++ {
		for y := tymin; y <= tymax; y++ {
			if x >= fxmin && x <= fxmax && y >= fymin && y <= fymax {
				continue
			}
			grid := &m.grids[x][y]
			grid.addWatcher(entity)
		}
	}
}

// 获取九宫格范围
func (m *AOIManager) getWatchGrids(pos Position) (int, int, int, int) {
	xmin, ymin := m.transXY(pos.x-m.gsize, pos.y-m.gsize)
	xmax, ymax := m.transXY(pos.x+m.gsize, pos.y+m.gsize)
	return xmin, xmax, ymin, ymax
}

// 遍历九宫格
func (m *AOIManager) visitWatchGrids(pos Position, f func(*Grid)) {
	xmin, xmax, ymin, ymax := m.getWatchGrids(pos)
	for x := xmin; x <= xmax; x++ {
		for y := ymin; y <= ymax; y++ {
			grid := &m.grids[x][y]
			f(grid)
		}
	}
}

// 获取坐标所在的网格
func (m *AOIManager) posToGrid(pos Position) *Grid {
	x, y := m.transXY(pos.x, pos.y)
	return &m.grids[x][y]
}

func (m *AOIManager) transXY(px, py float32) (int, int) {
	x := int((px - m.minX) / m.gsize)
	if x < 0 {
		x = 0
	}
	if x >= m.xNum {
		x = m.xNum - 1
	}

	y := int((py - m.minY) / m.gsize)
	if y < 0 {
		y = 0
	} else if y >= m.yNum {
		y = m.yNum - 1
	}

	return x, y
}
