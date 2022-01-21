package pacman

import (
	"ebitenGame/mygame3/assets"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"math/rand"
	"time"
)

const (
	ScreenWidth  = 712
	ScreenHeight = 1220
)

const (
	GameLoading gameState = iota
	GameStart
	GamePause
	GameOver
)

const OffsetY = CellSize * 10

type gameState int

type Game struct {
	state       gameState
	rand        *rand.Rand
	maze        *Maze
	data        *Data
	skinView    func(gameState, *Data) (*ebiten.Image, error)
	gridView    func(gameState, *Data) (*ebiten.Image, error)
	direction   direction
	powerTicker *time.Ticker

	audio *Audio
}

func NewGame() (*Game, error) {
	lAssets, err := assets.LoadAssets()
	if err != nil {
		return nil, err
	}
	mazeView, err := MazeView(lAssets.Walls)
	if err != nil {
		return nil, err
	}

	gridView, err := GridView(lAssets.Characters, lAssets.Powers, lAssets.ArcadeFont, mazeView)
	if err != nil {
		return nil, err
	}

	skinView, err := SkinView(lAssets.Skin, lAssets.Powers, lAssets.ArcadeFont)
	if err != nil {
		return nil, err
	}

	return &Game{
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		state:    GameLoading,
		gridView: gridView,
		skinView: skinView,
	}, nil
}

func (g *Game) Update() error {
	switch g.state {
	case GameLoading:
		if spaceReleased() {
			xcol := g.rand.Intn(Columns)
			numOfRows := MazeViewSize / CellSize
			g.data = NewData()
			g.maze = NewPopulatedMaze(32, g.rand)
			g.data.grid = g.maze.Get(0, numOfRows)
			g.data.active = make([][Columns]bool, numOfRows, numOfRows)
			g.data.pacman = Pacman{
				Position{
					cellX:     xcol,
					cellY:     0,
					posX:      float64((xcol * CellSize) + (CellSize / 2)),
					posY:      CellSize / 2,
					direction: North,
				},
			}
			g.direction = g.data.pacman.direction
			g.data.active[0][xcol] = true

			powers := make([]Power, 0)
			for i := 0; i < numOfRows; i += 4 {
				cellX := g.rand.Intn(Columns)
				cellY := g.rand.Intn(4) + i
				kind := Invincibility
				if (cellY-i)%2 == 0 {
					kind = Life
				}
				powers = append(powers, NewPower(cellX, cellY, kind))
			}
			g.data.powers = powers

			ghosts := make([]Ghost, 0)
			for i := 0; i < numOfRows; i += 2 {
				cellX := g.rand.Intn(Columns/2) + Columns/2
				if i%4 == 0 {
					cellX = g.rand.Intn(Columns / 2)

				}
				cellY := g.rand.Intn(2) + i
				kind := Ghost1
				if (cellY-i)%4 == 0 {
					kind = Ghost4
				} else if (cellY-i)%3 == 0 {
					kind = Ghost3
				} else if (cellY-i)%2 == 0 {
					kind = Ghost2
				}
				ghosts = append(ghosts, NewGhost(cellX, cellY, kind, getExit(
					g.data.grid[cellY][cellX])))
			}
			g.data.ghosts = ghosts
			//todo audio

			g.state = GameStart
		} else {
			g.data = nil
			g.maze = nil
			//todo audio
		}
	case GameStart:
		if spaceReleased() {
			g.state = GamePause
		} else if g.data.lifes < 1 {
			g.state = GameOver
		} else {
			g.running()
		}
	case GamePause:
		if spaceReleased() {
			g.state = GameStart
		}
	case GameOver:
		if spaceReleased() {
			g.state = GameLoading
			//todo audio
		} else {
			//todo audio
		}
	default:
		g.state = GameLoading
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	sview, err := g.skinView(g.state, g.data)
	if err != nil {
		fmt.Println(err)
	}

	gview, err := g.gridView(g.state, g.data)
	if err != nil {
		fmt.Println(err)
	}

	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Reset()
	screen.DrawImage(sview, ops)

	ops.GeoM.Reset()
	ops.GeoM.Translate(38, 162)
	screen.DrawImage(gview, ops)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) running() {
	numOfRows := MazeViewSize / CellSize
	if g.data.pacman.cellY == len(g.data.grid)-8 {
		g.maze.Compact(4)
		if (g.maze.Rows() - numOfRows) < 4 {
			g.maze.GrowBy(16)
		}

		g.data.grid = g.maze.Get(0, numOfRows)
		for i := 0; i < len(g.data.active); i++ {
			for j := 0; j < Columns; j++ {
				if i <= g.data.pacman.cellY {
					g.data.active[i-4][j] = g.data.active[i][j]
				} else {
					g.data.active[i-4][j] = false
				}
			}
		}

		g.data.pacman.cellY -= 4
		g.data.gridOffsetY -= CellSize * 4

		for i := 0; i < len(g.data.powers); i++ {
			g.data.powers[i].cellY -= 4
			if g.data.powers[i].cellY < 0 {
				cellX := g.rand.Intn(Columns)
				cellY := g.rand.Intn(4) + (numOfRows - 4)
				g.data.powers[i] = NewPower(cellX, cellY, g.data.powers[i].kind)
			}
		}

		for i := 0; i < len(g.data.ghosts); i++ {
			g.data.ghosts[i].cellY -= 4
			g.data.ghosts[i].posY -= CellSize * 4
			if g.data.ghosts[i].cellY < 0 {
				cellX := g.rand.Intn(Columns)
				cellY := g.rand.Intn(4) + (numOfRows - 4)
				g.data.powers[i] = NewPower(cellX, cellY, g.data.powers[i].kind)
			}
		}
		for i := 0; i < len(g.data.ghosts); i++ {
			g.data.ghosts[i].cellY -= 4
			g.data.ghosts[i].posY -= CellSize * 4
			if g.data.ghosts[i].cellY < 0 {
				cellX := g.rand.Intn(Columns)
				cellY := g.rand.Intn(4) + (numOfRows - 4)
				g.data.ghosts[i] = NewGhost(
					cellX, cellY,
					g.data.ghosts[i].kind,
					getExit(g.data.grid[cellY][cellX]))
			}
		}
	}

	g.keybord()
	g.movePacman()

	if !g.data.active[g.data.pacman.cellY][g.data.pacman.cellX] {
		if math.Abs(float64((g.data.pacman.cellX*CellSize)+(CellSize/2))-(g.data.pacman.posX)) < 20 &&
			math.Abs(float64((g.data.pacman.cellY*CellSize)+(CellSize/2))-(g.data.pacman.posY+g.data.gridOffsetY)) < 20 {
			g.data.active[g.data.pacman.cellY][g.data.pacman.cellX] = true
			g.data.score += 1
			//todo audio
		}
	}

	for i := 0; i < len(g.data.powers); i++ {
		cellX := g.rand.Intn(Columns)
		cellY := g.rand.Intn(4) + (((g.data.powers[i].cellY / 4) * 4) + numOfRows)
		if g.pacmanTouchesPower(i) {
			switch g.data.powers[i].kind {
			case Life:
				if g.data.lifes < MaxLifes {
					g.data.lifes += 1
					g.data.powers[i] = NewPower(cellX, cellY, g.data.powers[i].kind)
					//todo audio
				}
			case Invincibility:
				if !g.data.invincible {
					g.data.invincible = true
				}
				g.startCountdown(10)
				g.data.powers[i] = NewPower(cellX, cellY, g.data.powers[i].kind)
				//todo audio
			}
		}
	}
	for i := 0; i < len(g.data.ghosts); i++ {
		if g.pacmanTouchesGhost(i) {
			if !g.data.invincible {
				g.data.lifes -= 1
			} else {
				g.data.score += 200
				//todo audio
			}
			cellX := g.rand.Intn(Columns)
			cellY := g.rand.Intn(4) + (((g.data.ghosts[i].cellY / 4) * 4) + numOfRows)
			g.data.ghosts[i] = NewGhost(cellX, cellY, g.data.ghosts[i].kind, North)
		}
		g.moveGhost(i)
	}
}

func (g *Game) startCountdown(duration int) {
	if g.powerTicker != nil {
		g.powerTicker.Stop()
	}

	g.powerTicker = time.NewTicker(time.Duration(duration) * time.Second)
	go func() {
		select {
		case <-g.powerTicker.C:
			g.data.invincible = false
		}
	}()
}

func (g *Game) movePacman() {
	speed := 2.0
	xcell := g.data.pacman.cellX
	ycell := g.data.pacman.cellY

	switch g.direction {
	case North, South:
		if g.data.pacman.posX == float64((CellSize*xcell)+(CellSize/2)) {
			g.data.pacman.direction = g.direction
		}
	case East, West:
		if g.data.pacman.posY+g.data.gridOffsetY == float64((CellSize*ycell)+(CellSize/2)) {
			g.data.pacman.direction = g.direction
		}
	}

	switch g.data.pacman.direction {
	case North:
		if canMove(20.0,
			g.data.pacman.posX,
			g.data.pacman.posY+g.data.gridOffsetY+speed,
			g.data.pacman.cellY,
			g.data.pacman.cellY,
			g.data.grid[ycell][xcell],
		) {
			if g.data.pacman.posY > OffsetY {
				g.data.gridOffsetY += speed
			} else {
				g.data.pacman.posY += speed
			}
			if g.data.pacman.posY+g.data.gridOffsetY+20 > float64((ycell*CellSize)+CellSize) {
				g.data.pacman.cellY += 1
			}
		}
	case South:
		if canMove(20.0,
			g.data.pacman.posX,
			g.data.pacman.posY+g.data.gridOffsetY-speed,
			g.data.pacman.cellX,
			g.data.pacman.cellY,
			g.data.grid[ycell][xcell],
		) {
			if g.data.pacman.posY > OffsetY && g.data.gridOffsetY > 0 {
				g.data.gridOffsetY -= speed
			} else {
				g.data.pacman.posY -= speed
			}
			if g.data.pacman.posY+g.data.gridOffsetY-20 < float64(ycell*CellSize) {
				g.data.pacman.cellY -= 1
			}
		}
	case East:
		if canMove(20.0,
			g.data.pacman.posX+speed,
			g.data.pacman.posY+g.data.gridOffsetY,
			g.data.pacman.cellX,
			g.data.pacman.cellY,
			g.data.grid[ycell][xcell],
		) {
			g.data.pacman.posX += speed
			if g.data.pacman.posX+20 > float64((xcell*CellSize)+CellSize) {
				g.data.pacman.cellX += 1
			}
		}
	case West:
		if canMove(20.0,
			g.data.pacman.posX-speed,
			g.data.pacman.posY+g.data.gridOffsetY,
			g.data.pacman.cellX,
			g.data.pacman.cellY,
			g.data.grid[ycell][xcell],
		) {
			g.data.pacman.posX -= speed
			if g.data.pacman.posX-20 < float64(xcell*CellSize) {
				g.data.pacman.cellX -= 1
			}
		}
	}
}

func (g *Game) getGhostDirection(i int) direction {
	pacX := float64((g.data.pacman.cellX * CellSize) + (CellSize / 2))
	pacY := float64((g.data.pacman.cellY * CellSize) + (CellSize / 2))
	ghost := g.data.ghosts[i]

	if g.data.pacman.cellX == ghost.cellX && g.data.pacman.cellY == ghost.cellY {
		return ghost.direction
	}

	x, y := ghost.cellX, ghost.cellY
	ghsX := float64((x * CellSize) + (CellSize / 2))
	ghsY := float64((y * CellSize) + (CellSize / 2))

	prevDist := float64((MazeViewSize/CellSize)*Columns) * CellSize
	if g.data.invincible {
		prevDist = 0.0
	}

	for j := range g.rand.Perm(4) {
		if g.data.grid[ghost.cellY][ghost.cellX][j] == '_' {
			nx, ny := 0, 0
			dist := 0.0
			switch j {
			case 0: // North
				if y+1 < MazeViewSize/CellSize {
					dist = math.Sqrt(math.Pow(ghsX-pacX, 2) + math.Pow((ghsY+CellSize)-pacY, 2))
					nx, ny = ghost.cellX, ghost.cellY+1
				}
			case 1: // East
				dist = math.Sqrt(math.Pow((ghsX+CellSize)-pacX, 2) + math.Pow(ghsY-pacY, 2))
				nx, ny = ghost.cellX+1, ghost.cellY
			case 2: // South
				dist = math.Sqrt(math.Pow(ghsX-pacX, 2) + math.Pow((ghsY-CellSize)-pacY, 2))
				nx, ny = ghost.cellX, ghost.cellY-1
			case 3: // West
				dist = math.Sqrt(math.Pow((ghsX-CellSize)-pacX, 2) + math.Pow(ghsY-pacY, 2))
				nx, ny = ghost.cellX-1, ghost.cellY
			}
			if g.directionOfCell(ghost.cellX, ghost.cellY, nx, ny) !=
				getOppositeDirection(ghost.direction) {
				if g.data.invincible {
					if dist > prevDist {
						x, y, prevDist = nx, ny, dist
					}
				} else {
					if dist < prevDist {
						x, y, prevDist = nx, ny, dist
					}
				}
			}
		}
	}

	return g.directionOfCell(ghost.cellX, ghost.cellY, x, y)

}

func (g *Game) moveGhost(i int) {
	speed := 1.0
	ghost := g.data.ghosts[i]
	if ghost.cellY >= MazeViewSize/CellSize {
		return
	}
	if isIntersection(g.data.grid[ghost.cellY][ghost.cellX]) {
		if ghost.posX == float64((CellSize*ghost.cellX)+(CellSize/2)) &&
			ghost.posY == float64((CellSize*ghost.cellY)+(CellSize/2)) {
			g.data.ghosts[i].direction = g.getGhostDirection(i)
		}
	} else if isBlocked(g.data.grid[ghost.cellY][ghost.cellX], ghost.direction) || isDeadend(g.data.grid[ghost.cellY][ghost.cellX]) {
		if ghost.posX == float64((CellSize*ghost.cellX)+(CellSize/2)) &&
			ghost.posY == float64((CellSize*ghost.cellY)+(CellSize/2)) {
			g.data.ghosts[i].direction = getExit(g.data.grid[ghost.cellY][ghost.cellX])
		}
	}

	switch g.data.ghosts[i].direction {
	case North:
		if canMove(
			20.0,
			g.data.ghosts[i].posX,
			g.data.ghosts[i].posY+speed,
			g.data.ghosts[i].cellX,
			g.data.ghosts[i].cellY,
			g.data.grid[g.data.ghosts[i].cellY][g.data.ghosts[i].cellX],
		) {
			g.data.ghosts[i].posY += speed
			if g.data.ghosts[i].posY+20 > float64((g.data.ghosts[i].cellY*CellSize)+CellSize) {
				g.data.ghosts[i].cellY += 1
			}
		}
	case South:
		if canMove(
			20.0,
			g.data.ghosts[i].posX,
			g.data.ghosts[i].posY-speed,
			g.data.ghosts[i].cellX,
			g.data.ghosts[i].cellY,
			g.data.grid[g.data.ghosts[i].cellY][g.data.ghosts[i].cellX],
		) {
			g.data.ghosts[i].posY -= speed
			if g.data.ghosts[i].posY-20 < float64(g.data.ghosts[i].cellY*CellSize) {
				g.data.ghosts[i].cellY -= 1
			}
		}
	case East:
		if canMove(
			20.0,
			g.data.ghosts[i].posX+speed,
			g.data.ghosts[i].posY,
			g.data.ghosts[i].cellX,
			g.data.ghosts[i].cellY,
			g.data.grid[ghost.cellY][ghost.cellX],
		) {
			g.data.ghosts[i].posX += speed
			if g.data.ghosts[i].posX+20 > float64((g.data.ghosts[i].cellX*CellSize)+CellSize) {
				g.data.ghosts[i].cellX += 1
			}
		}
	case West:
		if canMove(
			20.0,
			g.data.ghosts[i].posX-speed,
			g.data.ghosts[i].posY,
			g.data.ghosts[i].cellX,
			g.data.ghosts[i].cellY,
			g.data.grid[g.data.ghosts[i].cellY][g.data.ghosts[i].cellX],
		) {
			g.data.ghosts[i].posX -= speed
			if g.data.ghosts[i].posX-20 < float64(g.data.ghosts[i].cellX*CellSize) {
				g.data.ghosts[i].cellX -= 1
			}
		}
	}
}

func (g *Game) pacmanTouchesPower(i int) bool {
	if g.data.powers[i].cellX == g.data.pacman.cellX && g.data.powers[i].cellY == g.data.pacman.cellY {
		posX := float64((g.data.powers[i].cellX * CellSize) + CellSize/2)
		posY := float64((g.data.powers[i].cellY * CellSize) + CellSize/2)
		if math.Abs(posX-g.data.pacman.posX) < 20 && math.Abs(posY-(g.data.pacman.posY+g.data.gridOffsetY)) < 20 {
			return true
		}
	}
	return false
}

func (g *Game) pacmanTouchesGhost(i int) bool {
	if g.data.ghosts[i].cellX == g.data.pacman.cellX && g.data.ghosts[i].cellY == g.data.pacman.cellY {
		posX := g.data.ghosts[i].posX
		posY := g.data.ghosts[i].posY
		if math.Abs(posX-g.data.pacman.posX) < 30 && math.Abs(posY-(g.data.pacman.posY+g.data.gridOffsetY)) < 30 {
			return true
		}
	}
	return false
}

func (g *Game) keybord() {
	if g.data != nil {
		walls := g.data.grid[g.data.pacman.cellY][g.data.pacman.cellX]
		if upKeyPressed() {
			if walls[0] == '_' {
				g.direction = North
			}
		}
		if downKeyPressed() {
			if walls[2] == '_' {
				g.direction = South
			}
		}
		if leftKeyPressed() {
			if walls[3] == '_' {
				g.direction = West
			}
		}
		if rightKeyPressed() {
			if walls[1] == '_' {
				g.direction = East
			}
		}
	}
}

func (g *Game) directionOfCell(cx, cy, nx, ny int) direction {
	if cx < nx {
		return East
	}
	if cx > nx {
		return West
	}
	if cy < ny {
		return North
	}
	if cy > ny {
		return South

	}
	if cx%2 == 0 {
		return West
	}
	return East
}

func getOppositeDirection(dir direction) direction {
	switch dir {
	case North:
		return South
	case East:
		return West
	case South:
		return North
	default:
		return East
	}
}

func numOfWalls(walls [4]rune) int {
	count := 0
	if walls[0] == 'N' {
		count += 1
	}
	if walls[1] == 'E' {
		count += 1
	}
	if walls[2] == 'S' {
		count += 1
	}
	if walls[3] == 'W' {
		count += 1
	}
	return count
}

func isIntersection(walls [4]rune) bool {
	count := numOfWalls(walls)
	if count >= 3 {
		return false
	} else if count == 2 {
		if walls[0] == walls[2] || walls[1] == walls[3] {
			return false
		}
	}
	return true
}

func isDeadend(walls [4]rune) bool {
	if numOfWalls(walls) >= 3 {
		return true
	}
	return false
}

func getExit(walls [4]rune) direction {
	for i := 0; i < 4; i++ {
		if walls[i] == '_' {
			switch i {
			case 0:
				return North
			case 1:
				return East
			case 2:
				return South
			case 3:
				return West
			}
		}
	}
	return North
}

func isBlocked(walls [4]rune, dir direction) bool {
	switch dir {
	case North:
		if walls[0] != '_' {
			return true
		}
	case East:
		if walls[1] != '_' {
			return true
		}
	case South:
		if walls[2] != '_' {
			return true
		}
	case West:
		if walls[3] != '_' {
			return true
		}
	}
	return false
}

func canMove(size float64, posX, posY float64, x, y int, walls [4]rune) bool {
	psx := posX - size
	psy := posY - size
	pex := posX + size
	pey := posY + size

	sx := x * CellSize
	sy := y * CellSize
	ex := sx + CellSize
	ey := sy + CellSize

	if walls[0] == 'N' {
		if pey > float64(ey-12) {
			return false
		}
	}
	if walls[1] == 'E' {
		if pex > float64(ex-12) {
			return false
		}
	}
	if walls[2] == 'S' {
		if psy > float64(sy+12) {
			return false
		}
	}
	if walls[3] == 'W' {
		if psx > float64(sx+12) {
			return false
		}
	}

	if pey > float64(ey-12) && psx < float64(sx+12) {
		return false
	}
	if pey > float64(ey-12) && pex > float64(ex-12) {
		return false
	}
	if psy < float64(sy+12) && psx < float64(sx+12) {
		return false
	}
	if psy < float64(sy+12) && pex > float64(ex-12) {
		return false
	}

	return true
}
