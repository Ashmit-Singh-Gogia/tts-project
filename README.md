# Full-Stack Text-to-Speech Engine

A full-stack web application that seamlessly converts text into streamable speech audio.

## Tech Stack

**Frontend:**
* React.js (Scaffolded with Vite for speed)
* Tailwind CSS (For modern, responsive UI styling)
* JavaScript Fetch API (For backend communication)

**Backend:**
* Go (Golang)
* Gin Web Framework (Routing & API endpoints)
* GORM with SQLite (Database management)
* `htgotts` (Native Google Text-to-Speech generation)

## 📂 Project Structure

The repository is divided into two main environments:
* `/backend` - The Go API, SQLite database, and audio generation engine.
* `/frontend` - The React user interface.
