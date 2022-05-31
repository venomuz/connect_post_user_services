package postgres

import (
	"github.com/jmoiron/sqlx"
	pb "github.com/venomuz/project3/PostService/genproto"
)

type postRepo struct {
	db *sqlx.DB
}

//NewPostRepo ...
func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) PostCreate(post *pb.Post) (*pb.OkBOOL, error) {
	PostQuery := `INSERT INTO posts(id,name,description,user_id) VALUES($1,$2,$3,$4)`

	_, err := r.db.Exec(PostQuery, post.Id, post.Name, post.Description, post.UserId)
	if err != nil {
		return nil, err
	}
	for _, media := range post.Medias {
		PostQuery := `INSERT INTO medias(id,post_id,type,link) VALUES($1,$2,$3,$4)`
		_, err := r.db.Exec(PostQuery, media.Id, post.Id, media.Type, media.Link)
		if err != nil {
			return nil, err
		}
	}

	return &pb.OkBOOL{Status: true}, nil
}
func (r *postRepo) PostGetByID(ID string) (*pb.Post, error) {
	post := pb.Post{}
	GetPostQuery := `SELECT id,name,description FROM posts WHERE user_id = $1`
	err := r.db.QueryRow(GetPostQuery, ID).Scan(&post.Id, &post.Name, &post.Description)
	if err != nil {
		return nil, err
	}
	var medias []*pb.Media
	GetMediaQuery := `SELECT id,post_id,type,link FROM medias WHERE post_id = $1`
	rows, err := r.db.Query(GetMediaQuery, post.Id)
	for rows.Next() {
		media := pb.Media{}
		err := rows.Scan(&media.Id, &media.PostId, &media.Type, &media.Link)
		if err != nil {
			return nil, err
		}
		medias = append(medias, &media)
	}
	post.Medias = medias

	return &post, nil
}
func (r *postRepo) PostDeleteByID(ID string) (*pb.OkBOOL, error) {
	_, err := r.db.Exec(`DELETE  FROM posts WHERE user_id = $1`, ID)
	if err != nil {
		return nil, err
	}

	return &pb.OkBOOL{Status: true}, nil
}
