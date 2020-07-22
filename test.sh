#!/bin/sh

echo "== bob registers"
curl -X POST localhost:3001/sessions/ -H 'Content-Type: application/json' -d '{"login": "bob", "password": "pass", "address": "bob:4000"}' -H 'user:bob'

echo "== alice registers"
curl -X POST localhost:3002/sessions/ -H 'Content-Type: application/json' -d '{"login": "alice", "password": "pass", "address": "alice:4000"}' -H 'user:alice'

echo "== alice sends a message to bob"
curl -X POST localhost:3002/messages/ -H 'Content-Type: application/json' -d '{"message":"salut bob, c est alice", "To": "bob"}' -H "user: alice"

echo "== bob replies"
curl -X POST localhost:3001/messages/ -H 'Content-Type: application/json' -d '{"message":"salut alice !", "To": "alice"}' -H "user: bob"

echo "== alice responds to bob"
curl -X POST localhost:3002/messages/ -H 'Content-Type: application/json' -d '{"message":"salut bob, c est alice", "To": "bob"}' -H "user: alice"

echo
echo "== alice checks messages from bob"
curl localhost:3002/conversations/bob

echo
echo
echo "== bob checks messages from alice"
curl localhost:3001/conversations/alice
