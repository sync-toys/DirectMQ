cmake_minimum_required(VERSION 3.26)

include(FetchContent)

# Fetch libwebsockets
FetchContent_Declare(
  libwebsockets
  GIT_REPOSITORY https://github.com/warmcat/libwebsockets.git
  GIT_TAG v4.3.3
)

set(LWS_WITHOUT_TESTAPPS ON CACHE BOOL "Disable libwebsockets test apps")
FetchContent_MakeAvailable(libwebsockets)

# Fetch nanopb
FetchContent_Declare(
  nanopb
  GIT_REPOSITORY https://github.com/nanopb/nanopb.git
  GIT_TAG 0.4.9
)
FetchContent_MakeAvailable(nanopb)

add_subdirectory(dmq_proto_embedded)
