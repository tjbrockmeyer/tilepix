package tilepix

import (
	"fmt"

	"github.com/faiface/pixel"
)

/*
  _____ _ _
 |_   _(_) |___
   | | | | / -_)
   |_| |_|_\___|
*/

// Tile is a TMX file structure which holds a Tiled tile.
type Tile struct {
	ID    ID     `xml:"id,attr"`
	Image *Image `xml:"image"`

	// parentMap is the map which contains this object
	parentMap *Map
}

func (t *Tile) String() string {
	return fmt.Sprintf("Tile{ID: %d}", t.ID)
}

func (t *Tile) setParent(m *Map) {
	t.parentMap = m

	if t.Image != nil {
		t.Image.setParent(m)
	}
}

// DecodedTile is a convenience struct, which stores the decoded data from a Tile.
type DecodedTile struct {
	ID             ID
	Tileset        *Tileset
	HorizontalFlip bool
	VerticalFlip   bool
	DiagonalFlip   bool
	Nil            bool

	sprite *pixel.Sprite
	pos    pixel.Vec

	// parentMap is the map which contains this object
	parentMap *Map
}

// Draw will draw the tile to the target provided.  This will calculate the sprite from the provided tileset and set the
// DecodedTiles' internal `sprite` property; this is so it is only calculated the first time.
func (t *DecodedTile) Draw(ind, columns, numRows int, ts *Tileset, target pixel.Target, offset pixel.Vec) {
	if t.IsNil() {
		return
	}

	if t.sprite == nil {
		t.setSprite(columns, numRows, ts)

		// Calculate the framing for the tile within its tileset's source image
		gamePos := indexToGamePos(ind, t.parentMap.Width, t.parentMap.Height)
		t.pos = gamePos.ScaledXY(pixel.V(float64(ts.TileWidth), float64(ts.TileHeight))).Add(pixel.V(float64(ts.TileWidth), float64(ts.TileHeight)).Scaled(0.5))
	}
	t.sprite.Draw(target, pixel.IM.Moved(t.pos.Add(offset)))
}

func (t *DecodedTile) String() string {
	return fmt.Sprintf("DecodedTile{ID: %d, Is nil: %t}", t.ID, t.Nil)
}

// IsNil returns whether this tile is nil.  If so, it means there is nothing set for the tile, and should be skipped in
// drawing.
func (t *DecodedTile) IsNil() bool {
	return t.Nil
}

func (t *DecodedTile) setParent(m *Map) {
	t.parentMap = m
}

func (t *DecodedTile) setSprite(columns, numRows int, ts *Tileset) {
	if t.IsNil() {
		return
	}

	if t.sprite == nil {
		// Calculate the framing for the tile within its tileset's source image
		x, y := tileIDToCoord(t.ID, columns, numRows)
		iX := float64(x)*float64(ts.TileWidth) + float64(ts.Margin+ts.Spacing*(x-1))
		fX := iX + float64(ts.TileWidth)
		iY := float64(y)*float64(ts.TileHeight) + float64(ts.Margin+ts.Spacing*(y-1))
		fY := iY + float64(ts.TileHeight)

		t.sprite = pixel.NewSprite(ts.sprite.Picture(), pixel.R(iX, iY, fX, fY))
	}
}
