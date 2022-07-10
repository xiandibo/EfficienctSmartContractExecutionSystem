prepare:
	@sudo rm -rf /tmp/badger
	@go build
	@cd examples/coordinatorclient && go build
	@mkdir /tmp/badger
	@mkdir /tmp/badger/coordinator
	@mkdir /tmp/badger/follower

run-example-coordinator:
	@./committer -role=coordinator -nodeaddr=10.20.82.234:3000 -followers=localhost:3001 -committype=two-phase -timeout=1000 -dbpath=/tmp/badger/coordinator -whitelist=127.0.0.1

run-example-follower1:
	@./committer -role=follower -nodeaddr=10.20.3.46:3001 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower2:
	@./committer -role=follower -nodeaddr=10.20.3.46:3002 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower3:
	@./committer -role=follower -nodeaddr=10.20.3.46:3003 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower4:
	@./committer -role=follower -nodeaddr=10.20.3.46:3004 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower5:
	@./committer -role=follower -nodeaddr=10.20.6.199:3001 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower6:
	@./committer -role=follower -nodeaddr=10.20.6.199:3002 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower7:
	@./committer -role=follower -nodeaddr=10.20.6.199:3003 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower8:
	@./committer -role=follower -nodeaddr=10.20.6.199:3004 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower9:
	@./committer -role=follower -nodeaddr=10.20.36.229:3001 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower10:
	@./committer -role=follower -nodeaddr=10.20.36.229:3002 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower11:
	@./committer -role=follower -nodeaddr=10.20.36.229:3003 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower12:
	@./committer -role=follower -nodeaddr=10.20.36.229:3004 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower13:
	@./committer -role=follower -nodeaddr=10.20.0.57:3001 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower14:
	@./committer -role=follower -nodeaddr=10.20.0.57:3002 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower15:
	@./committer -role=follower -nodeaddr=10.20.0.57:3003 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-follower16:
	@./committer -role=follower -nodeaddr=10.20.0.57:3004 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1




run-example-voter-4111:
	@./committer -role=follower -nodeaddr=10.20.6.199:4111 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4112:
	@./committer -role=follower -nodeaddr=10.20.36.229:4112 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4121:
	@./committer -role=follower -nodeaddr=10.20.0.57:4121 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4122:
	@./committer -role=follower -nodeaddr=10.20.6.199:4122 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4131:
	@./committer -role=follower -nodeaddr=10.20.36.229:4131 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4132:
	@./committer -role=follower -nodeaddr=10.20.0.57:4132 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1



run-example-voter-4211:
	@./committer -role=follower -nodeaddr=10.20.3.46:4211 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4212:
	@./committer -role=follower -nodeaddr=10.20.36.229:4212 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4221:
	@./committer -role=follower -nodeaddr=10.20.0.57:4221 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4222:
	@./committer -role=follower -nodeaddr=10.20.3.46:4222 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4231:
	@./committer -role=follower -nodeaddr=10.20.36.229:4231 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4232:
	@./committer -role=follower -nodeaddr=10.20.0.57:4232 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1



run-example-voter-4311:
	@./committer -role=follower -nodeaddr=10.20.3.46:4311 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4312:
	@./committer -role=follower -nodeaddr=10.20.6.199:4312 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4321:
	@./committer -role=follower -nodeaddr=10.20.0.57:4321 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4322:
	@./committer -role=follower -nodeaddr=10.20.3.46:4322 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4331:
	@./committer -role=follower -nodeaddr=10.20.6.199:4331 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4332:
	@./committer -role=follower -nodeaddr=10.20.0.57:4332 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1



run-example-voter-4411:
	@./committer -role=follower -nodeaddr=10.20.3.46:4411 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4412:
	@./committer -role=follower -nodeaddr=10.20.6.199:4412 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4421:
	@./committer -role=follower -nodeaddr=10.20.36.229:4421 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4422:
	@./committer -role=follower -nodeaddr=10.20.3.46:4422 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4431:
	@./committer -role=follower -nodeaddr=10.20.6.199:4431 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-voter-4432:
	@./committer -role=follower -nodeaddr=10.20.36.229:4432 -committype=two-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1



run-example-client:
	@examples/coordinatorclient/coordinator

unit-tests:
	@cd server && go test

functional-tests:
	@go test
