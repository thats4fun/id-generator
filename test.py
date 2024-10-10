# This utility:

# 1. Should be able to be called from a bash script/ the command line.
# 2. Should implement the following API:
# 3. Able to generate an ID that has not already been generated and output to the command line.
# 4. Able to free an already generated ID and reuse it the next time an ID is needed.
# 5. Should be concurrency-safe.

# Starter code:
# python
import argparse

# Please refer to README for more instructions
def getId():
    print("return a unique ID")

def freeId(id_to_free):
    print("free the id")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Manage unique integers')
    parser.add_argument("-d", "--deallocate",dest = "uniq_id_to_delete", help="Unique ID to delete", type=int)

    args = parser.parse_args()
    if args.uniq_id_to_delete != None:
        freeId(args.uniq_id_to_delete)
    else:
        uniq_id = getId()
    if uniq_id != None:
        print(uniq_id)
