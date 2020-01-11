# osMigrate

osMigrate (Open-source Migrate) was a project created due to [Flyway](https://flywaydb.org/)'s MySQL 5.6 support becoming an enterprise only
option. At work we have our database migrations running using Flyway and could not upgrade our MySQL version due to Amazon RDS features.
We also didn't fancy paying thousands of dollars for Flyway Enterprise, so osMigrate was born as a replacement for our current
Flyway configuration files.

It attempts to mimic the behavior of Flyway as close as possible, but is not a drop in replacement. At the moment it only
runs MySQL/MariaDB database migrations.