
set(NANOPB_INCLUDE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/nanopb")
include_directories(${NANOPB_INCLUDE_DIR})

add_subdirectory(nanopb)
add_subdirectory(dmq-proto-embedded)
