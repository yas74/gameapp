# Game Application

# use-case

## User use-cases

### Register
user can register to application by phone number

### Login
user can log in to application by phone number and password

## Game use-cases
### Each game has a given number of questions
### The dificulty levels of questions are "easy, medium, hard"
### Game winner is the one with most correct answers
### Each game falls into a specific category: sport, history, etc

# entity

## User
- ID
- Phone Number
- Avatar
- Name 

## Game
- ID
- Category
- Question List
- players
- Winer

## Question
- ID
- Question
- Answer List
- Correct Answer
- Difficulty
- Category

## Category
- ID
- Name
- Description