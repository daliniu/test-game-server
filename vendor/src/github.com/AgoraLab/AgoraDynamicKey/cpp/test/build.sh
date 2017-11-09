#!/bin/bash
g++ -ggdb -O0 DynamicKey_test.cpp DynamicKey2_test.cpp DynamicKey3_test.cpp DynamicKey4_test.cpp DynamicKey5_test.cpp main.cpp -lcrypto -std=c++0x -lgtest -lpthread -o TestDynamicKey.exe
