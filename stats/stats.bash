#!/bin/bash
#
census()
{
	echo -e "\tcensus: $1 $2 $3 $4"
	SWIZZLE_BTS=$1 SWIZZLE_NIAR=$2 SWIZZLE_REG=$3 ../bin/go build -a
	ropgadget --all --binary stats --depth 20 > "$4"
}

statistics()
{
	echo -e "statistics: $1 $2 $3 $4"

	rm $4-?.txt 2>/dev/null

	census $1 $2 $3 $4-1.txt
	census $1 $2 $3 $4-2.txt
	census $1 $2 $3 $4-3.txt
	census $1 $2 $3 $4-4.txt
	census $1 $2 $3 $4-5.txt
	
	cat $4-?.txt | grep '^0x' | sort | uniq -c | sort -nr | awk '{print $1}' | sort | uniq -c
	echo -e 

	rm $4-?.txt
}

statistics 0 0 0 NONE
statistics 50 0 0 BTS
statistics 0 50 0 NIAR
statistics 0 0 50 REG
statistics 50 50 50 ALL
