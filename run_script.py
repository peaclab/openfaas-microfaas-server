import sys
import json
src = json.loads(sys.argv[1])
params = json.loads(sys.argv[2])
exec(src)
handler(params)