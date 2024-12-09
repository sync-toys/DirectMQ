cmake_minimum_required(VERSION 3.26)

include(FetchContent)

# Fetch nanopb
FetchContent_Declare(
  nanopb
  GIT_REPOSITORY https://github.com/nanopb/nanopb.git
  GIT_TAG 0.4.9
)
FetchContent_MakeAvailable(nanopb)

# Fetch tinyhal
FetchContent_Declare(
  tinyhal
  GIT_REPOSITORY https://github.com/lexus2k/tinyhal
  GIT_TAG v0.1.1
)
FetchContent_MakeAvailable(tinyhal)

# Fetch tinyslip
FetchContent_Declare(
  tinyslip
  GIT_REPOSITORY https://github.com/lexus2k/tinyslip
  GIT_TAG v0.1.0
)
FetchContent_MakeAvailable(tinyslip)

add_subdirectory(dmq_proto_embedded)
