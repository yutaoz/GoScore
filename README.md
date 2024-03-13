# GoScore
Project to get historical basketball games from basketball-reference

## Details
Formatted with data folder with folders for each year, and text files named {month #}.txt. Each text file has a list of games, formatted as: Date, Home team, Away Team, Home Points, Away Points, Arena, Attendance, Box score url for bball ref.

Program currently limits 1 year per minute to avoid bball-refs load limiter. Can be changed to make ~20 requests per minute before getting session jailed, each year makes about 6 ish requests (for the months).

## Roadmap
- Get box score data for games also
- Handle upcoming games with no score, attendance, and box score url
