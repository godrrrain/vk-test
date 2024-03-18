package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Movie struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Release_date string `json:"release_date"`
	Rating       int    `json:"rating"`
}

type Actor struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}

type MovieTitle struct {
	Title string `json:"title"`
}

type ActorName struct {
	Name string `json:"name"`
}

type MovieInfo struct {
	ID           int         `json:"id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	Release_date string      `json:"release_date"`
	Rating       int         `json:"rating"`
	Actors       []ActorName `json:"actors"`
}

type ActorInfo struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	Gender   string       `json:"gender"`
	Birthday string       `json:"birthday"`
	Movies   []MovieTitle `json:"movies"`
}

type Storage interface {
	GetMovies(ctx context.Context, sortField string) ([]MovieInfo, error)
	GetActors(ctx context.Context) ([]ActorInfo, error)
	CreateMovie(ctx context.Context, title string, description string, release_date string, rating int, actors []int) error
	CreateActor(ctx context.Context, name, gender, birthday string) error
	DeleteMovie(ctx context.Context, id int) error
	DeleteActor(ctx context.Context, id int) error
	UpdateMovie(ctx context.Context, id int, title string, description string, release_date string, rating int) error
	UpdateActor(ctx context.Context, id int, name string, gender string, birthday string) error
}

type postgres struct {
	db *pgxpool.Pool
}

func NewPgStorage(ctx context.Context, connString string) (*postgres, error) {
	var pgInstance *postgres
	var pgOnce sync.Once
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			fmt.Printf("Unable to create connection pool: %v\n", err)
			return
		}

		pgInstance = &postgres{db}
	})

	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func (pg *postgres) GetMovies(ctx context.Context, sortField string) ([]MovieInfo, error) {
	query := fmt.Sprintf(`SELECT id, title, description, release_date, rating FROM movie ORDER BY %s`, sortField)

	rows, err := pg.db.Query(ctx, query)

	var movies []Movie
	var moviesInfo []MovieInfo

	if err != nil {
		return moviesInfo, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	movies, err = pgx.CollectRows(rows, pgx.RowToStructByName[Movie])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return moviesInfo, err
	}

	for _, v := range movies {
		query = fmt.Sprintf(`SELECT actor.name from actor, movie_actor
		WHERE movie_actor.movie_id = %d and actor.id = movie_actor.actor_id;`, v.ID)

		rows, err := pg.db.Query(ctx, query)

		var actorNames []ActorName
		var movieInfo MovieInfo

		if err != nil {
			return moviesInfo, fmt.Errorf("unable to query: %w", err)
		}
		defer rows.Close()

		actorNames, err = pgx.CollectRows(rows, pgx.RowToStructByName[ActorName])
		if err != nil {
			fmt.Printf("CollectRows error: %v", err)
			return moviesInfo, err
		}

		movieInfo.ID = v.ID
		movieInfo.Title = v.Title
		movieInfo.Description = v.Description
		movieInfo.Release_date = v.Release_date
		movieInfo.Rating = v.Rating
		movieInfo.Actors = actorNames

		moviesInfo = append(moviesInfo, movieInfo)
	}

	return moviesInfo, nil
}

func (pg *postgres) GetActors(ctx context.Context) ([]ActorInfo, error) {

	query := `SELECT id, name, gender, birthday FROM actor`

	rows, err := pg.db.Query(ctx, query)

	var actors []Actor
	var actorsInfo []ActorInfo

	if err != nil {
		return actorsInfo, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	actors, err = pgx.CollectRows(rows, pgx.RowToStructByName[Actor])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return actorsInfo, err
	}

	for _, v := range actors {
		query = fmt.Sprintf(`SELECT movie.title from movie, movie_actor
		WHERE movie_actor.actor_id = %d and movie.id = movie_actor.movie_id;`, v.ID)

		rows, err := pg.db.Query(ctx, query)

		var movieTitles []MovieTitle
		var actorInfo ActorInfo

		if err != nil {
			return actorsInfo, fmt.Errorf("unable to query: %w", err)
		}
		defer rows.Close()

		movieTitles, err = pgx.CollectRows(rows, pgx.RowToStructByName[MovieTitle])
		if err != nil {
			fmt.Printf("CollectRows error: %v", err)
			return actorsInfo, err
		}

		actorInfo.ID = v.ID
		actorInfo.Name = v.Name
		actorInfo.Gender = v.Gender
		actorInfo.Birthday = v.Birthday
		actorInfo.Movies = movieTitles

		actorsInfo = append(actorsInfo, actorInfo)
	}

	return actorsInfo, nil
}

func (pg *postgres) CreateMovie(ctx context.Context, title string, description string, release_date string, rating int, actors []int) error {

	// query := `INSERT INTO movie (title, description, release_date, rating)
	// VALUES (@title, @description, @release_date, @rating)`
	// args := pgx.NamedArgs{
	// 	"title":        title,
	// 	"description":  description,
	// 	"release_date": release_date,
	// 	"rating":       rating,
	// }
	// _, err := pg.db.Exec(ctx, query, args)
	// if err != nil {
	// 	return fmt.Errorf("unable to insert row: %w", err)
	// }

	var id int
	err := pg.db.QueryRow(context.Background(), `INSERT INTO movie (title, description, release_date, rating) 
	VALUES ($1, $2, $3, $4) RETURNING id`, title, description, release_date, rating).Scan(&id)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	for _, v := range actors {
		query := `INSERT INTO movie_actor (movie_id, actor_id)
		VALUES (@movie_id, @actor_id)`
		args := pgx.NamedArgs{
			"movie_id": id,
			"actor_id": v,
		}
		_, err := pg.db.Exec(ctx, query, args)
		if err != nil {
			return fmt.Errorf("unable to insert row: %w", err)
		}
	}

	return nil
}

func (pg *postgres) CreateActor(ctx context.Context, name, gender, birthday string) error {

	query := `INSERT INTO actor (name, gender, birthday) 
	VALUES (@name, @gender, @birthday)`
	args := pgx.NamedArgs{
		"name":     name,
		"gender":   gender,
		"birthday": birthday,
	}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func (pg *postgres) DeleteMovie(ctx context.Context, id int) error {

	query := fmt.Sprintf(`DELETE FROM movie WHERE id = %d`, id)
	_, err := pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to delete row: %w", err)
	}

	query = fmt.Sprintf(`DELETE FROM movie_actor WHERE movie_id = %d`, id)
	_, err = pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to delete row: %w", err)
	}

	return nil
}

func (pg *postgres) DeleteActor(ctx context.Context, id int) error {

	query := fmt.Sprintf(`DELETE FROM actor WHERE id = %d`, id)
	_, err := pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to delete row: %w", err)
	}

	query = fmt.Sprintf(`DELETE FROM movie_actor WHERE actor_id = %d`, id)
	_, err = pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to delete row: %w", err)
	}

	return nil
}

func (pg *postgres) UpdateMovie(ctx context.Context, id int, title string, description string, release_date string, rating int) error {

	updateData := ""

	if title != "" {
		updateData += fmt.Sprintf(`title = '%s', `, title)
	}
	if description != "" {
		updateData += fmt.Sprintf(`description = '%s', `, description)
	}
	if release_date != "" {
		updateData += fmt.Sprintf(`release_date = '%s', `, release_date)
	}
	if rating != 0 {
		updateData += fmt.Sprintf(`rating = %d, `, rating)
	}
	if len(updateData) < 2 {
		return fmt.Errorf("fields to change must be specified")
	}
	updateData = updateData[:len(updateData)-2]

	query := fmt.Sprintf(`UPDATE movie SET %s WHERE id = %d`, updateData, id)

	_, err := pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to update row: %w", err)
	}

	return nil
}

func (pg *postgres) UpdateActor(ctx context.Context, id int, name string, gender string, birthday string) error {

	updateData := ""

	if name != "" {
		updateData += fmt.Sprintf(`name = '%s', `, name)
	}
	if gender != "" {
		updateData += fmt.Sprintf(`gender = '%s', `, gender)
	}
	if birthday != "" {
		updateData += fmt.Sprintf(`birthday = '%s', `, birthday)
	}
	if len(updateData) < 2 {
		return fmt.Errorf("fields to change must be specified")
	}
	updateData = updateData[:len(updateData)-2]

	query := fmt.Sprintf(`UPDATE actor SET %s WHERE id = %d`, updateData, id)

	_, err := pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to update row: %w", err)
	}

	return nil
}
