Remember to use ```docker compose up --build``` after making changes to code.

```docker compose --env-file ./dev.env up --build```

```docker compose down -v``` will remove volumes, so use it when you make a DB change.