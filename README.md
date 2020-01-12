# osMigrate

osMigrate (Open-source Migrate) was a project created due to [Flyway](https://flywaydb.org/)'s MySQL 5.6 support becoming an enterprise only
option. At work we have our database migrations running using Flyway and could not upgrade our MySQL version due to Amazon RDS features.
We also don't fancy paying thousands of dollars for Flyway Enterprise, so osMigrate was born as a replacement for our current
Flyway configuration files.

It attempts to mimic the behavior of Flyway as close as possible, but is not a drop in replacement. At the moment it only
runs MySQL/MariaDB database migrations.

# State of development

At the moment I would not use this in production. It is still in development and may not work properly and does not
have many of the safety features that Flyway comes with.

Additionally, this is my first Golang project, so I am still learning the language.

*I cannot be held responsible for any damage caused by running this software*

Contributions welcome.