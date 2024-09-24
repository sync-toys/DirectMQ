from conan import ConanFile
from conan.tools.cmake import cmake_layout

class DirectMQTestAgent(ConanFile):
    settings = "os", "compiler", "build_type", "arch"
    generators = ("CMakeToolchain", "CMakeDeps")

    def requirements(self):
        self.requires("nlohmann_json/3.11.3")

    def build_requirements(self):
        self.build_requires("cmake/[>3.25]")

    def layout(self):
        cmake_layout(self)
