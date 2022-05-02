import json

from privateChainTestTestSmartContract import compiled_sol_kvtest
from privateChainTestTestSmartContract import compiled_sol_kvstore


kvtest_out = json.loads(compiled_sol_kvtest['contracts']['kvtest.sol']['Kvs']['metadata'])

kvstore_out = json.loads(compiled_sol_kvstore['contracts']['kvstore.sol']['kvstore']['metadata'])

print("kvtest")
print(compiled_sol_kvtest['contracts']['kvtest.sol']['Kvs']['evm']['bytecode']['object'])
print("kvstore")
print(compiled_sol_kvstore['contracts']['kvstore.sol']['kvstore']['evm']['bytecode']['object'])

