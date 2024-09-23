from conan import ConanFile
from conan.tools.cmake import cmake_layout

class DirectMQ(ConanFile):
    settings = "os", "compiler", "build_type", "arch"
    generators = ("CMakeToolchain", "CMakeDeps")

    def requirements(self):
        self.requires("websocketpp/0.8.2")
        self.requires("uwebsockets/20.65.0")

    def build_requirements(self):
        self.build_requires("cmake/[>3.25]")

    def layout(self):
        cmake_layout(self)
