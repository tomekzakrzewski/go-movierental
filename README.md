FIX RENTING MOVIE
path -> apiv1/movies/:id/rent
struct -> time.Now - dodawanie czasu kiedy sie wypozyczylo now i +24 toDate
movieID -> z path
userID -> z contextu

--- ADD ---
sorting movies by rating, alphabetical, year
pagination
user can see their rented movies


USER -> auth, get movies, get movie by id, sort movie by rating, get rented movies by user, rate movie
CHECK IF MOVIE IS NOT ALREADY RENTED BY THE GIVERN USER
