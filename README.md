# gemini-dca

Personal automation scheduled logic to perform DCA into Gemini CEX

Refactored https://github.com/JeraldYik/gemini-dca into Golang implementation

[Link](https://docs.google.com/document/d/1jUYHuTD6vl7BIbh2C48RT9o6VYQupZ8BEPhZ4oCTtVs/edit?usp=sharing) to Design Document

## Setup

- [ ] Populate `conf/dev.env`/`conf/staging.env`/`conf/production.env` (depending on environment)
- [ ] Run docker and run command `docker-compose up`
- (If you want to remove the entire database functional group, comment and remove the relevant lines of code/files)

Refer to Makefile for executable commands

Note: If you are hosting on Heroku, it may be helpful to run `heroku scale web=0 --remote <remote-env>` to prevent `npm start` to be called on every deploy.
