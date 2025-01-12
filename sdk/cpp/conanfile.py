from conan import ConanFile
from conan.tools.cmake import cmake_layout

class DirectMQ(ConanFile):
    settings = "os", "compiler", "build_type", "arch"
    generators = ("CMakeToolchain", "CMakeDeps")

    def requirements(self):
        self.requires("boost/1.86.0")
        self.requires("catch2/2.13.10")

    def build_requirements(self):
        self.build_requires("cmake/[>3.25]")

    def layout(self):
        cmake_layout(self)
