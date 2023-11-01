#!/bin/bash

# Declare an array of filenames
file_names=(
    "englishclient_frontend.bsp.pak000_dir.vpk"
    "englishclient_mp_angel_city.bsp.pak000_dir.vpk"
    "englishclient_mp_black_water_canal.bsp.pak000_dir.vpk"
    "englishclient_mp_box.bsp.pak000_dir.vpk"
    "englishclient_mp_coliseum.bsp.pak000_dir.vpk"
)

# Loop through the array and create files
for name in "${file_names[@]}"; do
    # Generate a random size between 100 and 200 (in MB)
    size=$(( 100 + RANDOM % 101 ))

    # Use the dd command to create a file of the specified size
    dd if=/dev/urandom of="$name" bs=1M count=$size

    echo "Created $name with size $size MB"
done
