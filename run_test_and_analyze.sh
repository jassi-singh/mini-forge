#!/bin/bash

# Run k6 load test and analyze for duplicate keys
echo "🚀 Starting k6 load test..."
echo ""

# Run k6 and capture output
k6 run k6_load_test.js 2>&1 | tee test_output.log

echo ""
echo "📊 Extracting keys and analyzing for duplicates..."
echo ""

# Extract keys from the output (lines starting with "KEY:")
grep "KEY:" test_output.log | sed 's/.*KEY://' > results.log

# Count total keys
TOTAL_KEYS=$(wc -l < results.log | tr -d ' ')

# Count unique keys
UNIQUE_KEYS=$(sort results.log | uniq | wc -l | tr -d ' ')

# Find duplicates
DUPLICATE_COUNT=$(sort results.log | uniq -d | wc -l | tr -d ' ')

echo "============================================================"
echo "          KEY UNIQUENESS TEST RESULTS"
echo "============================================================"
echo ""
echo "Total Keys Generated:  $TOTAL_KEYS"
echo "Unique Keys:           $UNIQUE_KEYS"
echo "Duplicate Keys Found:  $DUPLICATE_COUNT"
echo ""

if [ "$DUPLICATE_COUNT" -eq 0 ]; then
    echo "✅ SUCCESS: All keys are unique!"
else
    echo "❌ FAILURE: Duplicate keys detected!"
    echo ""
    echo "Duplicate Keys Details:"
    echo "------------------------------------------------------------"
    # Show duplicate keys with their counts
    sort results.log | uniq -c | awk '$1 > 1 {print "  Key: " $2 " - Occurred " $1 " times"}' | head -20
    
    if [ "$DUPLICATE_COUNT" -gt 20 ]; then
        echo "  ... and more duplicates (check results.log for full list)"
    fi
fi

echo ""
echo "============================================================"
echo ""
echo "📁 All keys saved to: results.log"
echo "📁 Full test output saved to: test_output.log"
echo ""

