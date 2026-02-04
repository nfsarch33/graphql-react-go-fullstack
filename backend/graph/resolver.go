package graph

import (
	"gorm.io/gorm"
)

// This file will be automatically augmented by gqlgen when you run the generate command

// Resolver holds database connection for all resolvers
type Resolver struct {
	DB *gorm.DB
}
