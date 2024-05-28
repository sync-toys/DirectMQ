cmake_minimum_required(VERSION 3.26)

project(directmq-sdk
  VERSION 1.0
  LANGUAGES CXX)

find_package(Boost REQUIRED COMPONENTS system)

add_executable(directmq-sdk src/main.cpp)
target_link_libraries(directmq-sdk
  PRIVATE
  Boost::system)


option(TESTING "Enable tests." ON)

if(TESTING)
  enable_testing()
  add_subdirectory(test)
endif()

install(TARGETS directmq-sdk DESTINATION bin)