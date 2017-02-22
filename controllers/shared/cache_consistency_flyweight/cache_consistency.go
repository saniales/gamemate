//Package isCacheConsistent contains variables to check whenever cache is not consistent
//with the archives (so it should not be used and must be fixed from a routine).
//NOTE: cache is consistent when is empty or with ALL data from archives object.
//
//It is not consistent when even a single information in the archives is not
//correctly cached.
package isCacheConsistent

//User represents if the users' cache is consistent (sessions).
var User bool

//Developer represents if the developers' cache is consistent (sessions, api tokens).
var Developer bool

//GameOwner represents if the game owners' cache is consistent (sessions, owned games).
var GameOwner bool

//Match represents if the matches' cache is consistent (open matches).
var Match bool

//Game represents if the games' cache is consistent (game list, game summaries,
//game enabled for users).
var Game bool
