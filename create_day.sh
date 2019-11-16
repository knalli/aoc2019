#!/bin/bash

function createDay() {
  local day=$1
  local dayFmt
  local dayDir
  dayFmt=$(printf "%02d" "$day")
  dayDir="day${dayFmt}"
  if [ -d "day${dayDir}" ]; then
    echo "Day already exist"
  else
    echo "Creating Day $day in '${dayDir}/'..."
    cp -pR tpl "${dayDir}"
    if [[ "$OSTYPE" == "darwin"* ]] || [[ "$OSTYPE" == "freebsd"* ]]; then
      sed -i '' 's/AocDay = 1/AocDay = '"$day"'/g' "${dayDir}/main.go"
      sed -i '' 's/AocDayName = "day01"/AocDayName = "day'"$dayFmt"'"/g' "${dayDir}/main.go"
      sed -i '' 's/AocDayTitle = "Day 01"/AocDayTitle = "Day '"$dayFmt"'"/g' "${dayDir}/main.go"
    else
      sed -i 's/AocDay = 1/AocDay = '"$day"'/g' "${dayDir}/main.go"
      sed -i 's/AocDayName = "day01"/AocDayName = "day'"$dayFmt"'"/g' "${dayDir}/main.go"
      sed -i 's/AocDayTitle = "Day 01"/AocDayTitle = "Day '"$dayFmt"'"/g' "${dayDir}/main.go"
    fi
  fi
}

createDay "$1"
