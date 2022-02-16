go build
for file in ./test_scripts/*
do
    echo "Running tests in $file"
    ./becca $file
done