package db

import (
	"github.com/go-pg/pg"
	"github.com/polisgo2020/search-bolotrush/index"
)

type Token struct {
	ID   int    `sql:"id,pk"`
	Word string `sql:"word"`
	//Position int
}
type File struct {
	ID   int    `sql:"id,pk"`
	Name string `sql:"name"`
}
type Occurrence struct {
	ID       int `sql:"id,pk"`
	WordID   int `sql:"word_id"`
	FileID   int `sql:"file_id"`
	Position int `sql:"position"`
}
type Base struct {
	pg *pg.DB
}

func NewDb(config string) (*Base, error) {
	options, err := pg.ParseURL(config)
	if err != nil {
		return &Base{}, err
	}
	return &Base{
		pg: pg.Connect(options),
	}, nil
}

func (b *Base) WriteIndex(index index.InvMap) error {
	_, err := b.pg.Exec("TRUNCATE tokens, files, occurrences;")
	if err != nil {
		return err
	}
	err = b.addIndex(index)
	if err != nil {
		return err
	}
	return nil
}

func (b *Base) addIndex(index index.InvMap) error {
	for token, properties := range index {
		token := Token{
			Word: token,
		}
		err := b.pg.Insert(&token)
		if err != nil {
			return err
		}
		for _, property := range properties {
			file := File{
				Name: property.FileName,
			}
			err := b.pg.Insert(&file)
			if err != nil {
				return err
			}
			for _, position := range property.Positions {
				occurrence := Occurrence{
					WordID:   token.ID,
					FileID:   file.ID,
					Position: position,
				}
				err := b.pg.Insert(&occurrence)
				if err != nil {
					return err
				}

			}

		}
	}
	return nil
}

//
//func (b *Base) getToken(word string) (*Token, error) {
//
//}
//
//func (b *Base) GetMatches(rawQuery string) ([]index.MatchList, error) {
//	query := index.PrepareText(rawQuery)
//	if len(query) == 0 {
//		return nil, errors.New("wrong query")
//	}
//	var result []index.MatchList
//	word, err := b.getToken()
//}
