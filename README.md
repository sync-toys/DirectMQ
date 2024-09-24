# DirectMQ

- [What is DirectMQ?](#what-is-directmq)
- [Why DirectMQ?](#why-directmq)
- [When Not to Use DirectMQ](#when-not-to-use-directmq)
- [Key features](#key-features)
  * [1. Customizable Transport Layer](#1-customizable-transport-layer)
  * [2. Topic-Based Messaging](#2-topic-based-messaging)
  * [3. Wildcard Support for Topics](#3-wildcard-support-for-topics)
  * [4. Message Delivery Strategies](#4-message-delivery-strategies)
  * [5. Subscription Optimization](#5-subscription-optimization)
  * [6. Lightweight Protocol](#6-lightweight-protocol)
  * [7. Flexible Versioning](#7-flexible-versioning)
  * [8. Protocol Buffers Encoding](#8-protocol-buffers-encoding)
- [Getting Started](#getting-started)
- [Client Libraries](#client-libraries)
  * [Golang SDK](#golang-sdk)
  * [C++ SDK](#c-sdk)
  * [Dart SDK (Planned)](#dart-sdk-planned)
- [How It Works](#how-it-works)
  * [1. Protocol Version Negotiation](#1-protocol-version-negotiation)
  * [2. Connection Initialization](#2-connection-initialization)
  * [3. Subscription Synchronization](#3-subscription-synchronization)
  * [4. Message Exchange](#4-message-exchange)
  * [5. Graceful Disconnection](#5-graceful-disconnection)
- [Protocol Specification](#protocol-specification)
  * [Dockerized Tests Runner](#dockerized-tests-runner)
  * [Debugging Tools](#debugging-tools)
  * [Message Recording](#message-recording)
  * [Test Bench](#test-bench)
- [Development Environment](#development-environment)
  * [Cloud and Local Usage](#cloud-and-local-usage)
  * [Supported Code Editor](#supported-code-editor)
  * [Devcontainer Features](#devcontainer-features)
  * [Task Management](#task-management)
  * [Initial Setup](#initial-setup)
  * [Running Tests](#running-tests)
- [License](#license)
  * [Summary of the MIT License](#summary-of-the-mit-license)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

## What is DirectMQ?

**DirectMQ** is a lightweight, customizable messaging solution designed for peer-to-peer (P2P) communication between devices. It allows developers to send and receive messages on specified topics directly between programs without relying on predefined network layers. 

DirectMQ offers the flexibility to implement custom transport mechanisms, such as Bluetooth, Serial, or even radio communication, making it ideal for embedded systems and IoT devices.

## Why DirectMQ?

The need for a highly flexible P2P messaging system arose from the limitations of existing solutions like **nng** and **ZeroMQ**, which are often too network-dependent and not suited for low-resource or embedded platforms. These technologies do not support custom transport methods such as **Bluetooth** or **Serial** communication, which are essential for connecting smartphones to embedded devices.

With DirectMQ, you can define your own transport layer, giving you full control over how devices communicate. This makes it perfect for projects where traditional networking solutions are too restrictive or complex to implement.

## When Not to Use DirectMQ

DirectMQ is optimized for direct communication, especially between mobile applications and embedded/IoT devices, as well as between multiple embedded devices. However, there are situations where it is better to use other, proven communication protocols. Here are a few:

- **Complex System Architectures**: If your application requires a complex architecture with distributed components or needs advanced message management features, it is advisable to choose protocols like **MQTT**, **nng**, or **ZeroMQ**, which offer richer capabilities and have been well-tested in various conditions

- **High Reliability and Delivery Guarantees**: In scenarios where reliable message delivery and strict management are critical (e.g., in mission-critical applications), it is recommended to use protocols designed with these requirements in mind.

- **Scalability**: If you plan to scale your system to work with hundreds or thousands of devices, protocols like MQTT can provide better scalability and session management. This project is not meant to handle more than tens of devices connected in direct network.

- **Support for Industry Standards**: In projects that must meet specific industry standards, utilizing popular protocols with wide community support may be more appropriate.

- **Backend-to-Backend and Mobile-to-Mobile Communication**: DirectMQ is not recommended for backend-to-backend or mobile-to-mobile communication due to better alternatives available. Established protocols like **zeromq** are more suitable for these scenarios, providing robust features and better performance.

If you feel that MQTT might be a better choice but still want to use DirectMQ, it’s advisable to reconsider your architecture.

It is important to note that DirectMQ should not replace protocols like MQTT and _should only be used in contexts where it makes practical sense_. While DirectMQ is ideal for simple and direct connections, it is worth considering proven and established communication protocols in more complex scenarios.

## Key features

DirectMQ is designed to be a simple yet powerful messaging protocol that allows flexibility in communication between devices. Here are some of its key features:

### 1. Customizable Transport Layer

DirectMQ doesn't impose a predefined transport mechanism, allowing developers to implement custom transports such as:
- **Bluetooth**
- **Serial**
- **I2C**
- **WebSocket**
- **Radio communication**
- **Print on paper and send with pigeon 😜**

This flexibility makes DirectMQ ideal for embedded systems or specialized communication needs.

### 2. Topic-Based Messaging

Messages are exchanged based on user-defined topics, making it easy to organize and control communication flows. You can create hierarchical topics like:
- `a/b/c`
- `toy/some_toy_id/logs`
- `toy/some_toy_id/metrics`
- `toy/some_toy_id/connect_wifi`

### 3. Wildcard Support for Topics

DirectMQ supports wildcard topics, enabling flexible message subscriptions:
- Single-level wildcards: `a/*/c`
- Multi-level wildcards: `a/**`
- Multiple wildcards and mixing is also allowed

### 4. Message Delivery Strategies

DirectMQ provides two message delivery strategies to suit various use cases:
- **AT_LEAST_ONCE**: Ensures that messages are delivered at least once, good for pub/sub patterns.
- **AT_MOST_ONCE**: Ensures that messages are delivered at most once, without duplicates.

### 5. Subscription Optimization

DirectMQ optimizes network traffic by only transmitting messages for topics that nodes in the network have explicitly subscribed to, reducing unnecessary communication overhead.

The hierarchy of topics is also preserved, what means that if node on the left has subscribed to topic `a/b/*` and node in the middle subscribed to `a/**` then only subscription that is sent to the right node is `a/**` - no duplication of topics.

### 6. Lightweight Protocol

The protocol is kept simple, which minimizes resource usage. It does not include built-in delivery reports or health checks, making it lightweight and efficient for resource-constrained environments like embedded devices.

### 7. Flexible Versioning

DirectMQ supports version negotiation, allowing nodes to communicate using the highest supported version of the protocol. This ensures compatibility when updating the protocol in the future.

### 8. Protocol Buffers Encoding

Messages are encoded and decoded using **Protocol Buffers** (Protobuf), providing efficient binary encoding to minimize message size and improve performance. 

`C++ SDK` uses `nanopb` as protocol encoder/decoder ensuring that communication has minimal memory footprint.

## Getting Started

TODO

## Client Libraries

DirectMQ is built with flexibility in mind, allowing it to communicate across various devices and platforms using different programming languages. Below are the available and planned client libraries:

### Golang SDK

The **Golang SDK** was developed to create developer tools for the Synctoys project. It utilizes the standard Protocol Buffers generator and includes unit tests to ensure functionality and stability.

- **Use Case**: Ideal for writing tools for your mobile/embedded development environment.
- **Key Features**:
  - Built-in Protocol Buffers support.
  - Easy integration with existing Go projects.
- **Usage Example**:
    ```go
    // TODO: Golang SDK example usage
    ```

> *Status*: Fully implemented and tested.

### C++ SDK

The **C++ SDK** is optimized for embedded systems, where resource constraints like memory and processing power are critical. This SDK uses **nanopb**, a lightweight Protocol Buffers generator, ensuring minimal memory usage while maintaining performance.

- **Use Case**: Suitable for embedded systems, IoT devices, or any project where memory efficiency is a priority.
- **Key Features**:
  - Nanopb-generated Protocol Buffers for low memory footprint.
  - Optimized for use on devices with limited resources.
  - **Usage Example**:
    ```cpp
    // TODO: C++ SDK example usage
    ```

> *Status*: Implemented but not tested yet, missing published `platformio` and `conan index` registries packages.

### Dart SDK (Planned)

The **Dart SDK** is planned for applications built with the **Flutter** framework, focusing on mobile app development. This SDK will enable seamless communication between mobile apps and embedded devices running DirectMQ.

- **Use Case**: Perfect for developing mobile applications that need to communicate directly with IoT or embedded devices.
- **Key Features**:
  - Tight integration with Flutter apps.
  - Full compatibility with other DirectMQ SDKs for cross-platform communication.
  - **Usage Example**:
    ```dart
    // TODO: Dart SDK example usage
    ```

> *Status*: In development (Coming soon).

---

Each SDK is designed to work within the same DirectMQ ecosystem, ensuring smooth interoperability between different devices and programming languages.

## How It Works

DirectMQ operates through a series of steps to facilitate seamless communication between nodes. This section outlines the main processes involved in establishing a connection and exchanging messages.

### 1. Protocol Version Negotiation

The communication begins with nodes negotiating the protocol version they will use. This is done by sending a **Supported Protocol Versions** message, allowing both nodes to identify the highest common version available. Currently, the only available version is version 1. This ensures that both parties are aligned on the protocol specifications before any further communication takes place.

### 2. Connection Initialization

Once the protocol version is agreed upon, the next step involves initializing the connection between the nodes. At this stage, one of the nodes can choose to reject the connection request if necessary. This provides flexibility for nodes to control which connections they accept, ensuring that they only communicate with trusted partners.

### 3. Subscription Synchronization

After a successful connection initialization, the nodes synchronize their subscriptions. This involves exchanging information about active subscriptions, ensuring that both nodes are aware of which topics they are interested in. This synchronization is critical for establishing a common understanding of the message flow between the nodes.

### 4. Message Exchange

Once the subscription synchronization is complete, normal communication can begin. At this point, messages can be exchanged between the nodes based on the topics they have subscribed to. DirectMQ ensures that only messages that match the subscribed topics are transmitted, optimizing network usage, reducing unnecessary CPU and memory usage of nodes.

### 5. Graceful Disconnection

Once the communication session is complete, nodes can gracefully disconnect from each other. After disconnection, each node automatically optimizes its subscriptions with the remaining connected nodes. This ensures that the nodes maintain an efficient message flow by updating their subscription lists to reflect only the active connections. This process minimizes resource usage and enhances the overall performance of the network.

In case of forceful disconnection (like unplug of USB cable in case of SLIP transport) nodes try to graceful disconnect from each other, but such messages can be ignored by transport layer then, and normal edge removal can be performed. 

---

By following these steps, DirectMQ ensures a robust and flexible messaging framework that can be tailored to the specific needs of various applications, particularly in environments requiring direct peer-to-peer communication.

## Protocol Specification

The DirectMQ specification is written in the form of integration tests. Each client SDK runs an agent that tests its behaviors and communication through a series of common integration tests managed by test suite written in Go. This approach ensures that all SDKs adhere to the same protocol specifications, maintaining consistency across different programming environments.

### Dockerized Tests Runner

To maintain high stability independent of the host system, integration tests are run inside **Docker container** within a development container. This isolates the testing environment, ensuring that host system did not interfere with running tests.

### Debugging Tools

By simply toggling switches in the code, developers can activate debugging modes for each agent in the selected test. Visual Studio Code (VSCode) comes equipped with predefined debugging configurations, making it easier for maintainers to troubleshoot and monitor the behavior of their SDKs during testing.

### Message Recording

The specification tests record communication between the nodes being tested through a small WebSocket bridge and saves it as simple JSON snapshot files. Reviewing these communication logs allows developers to understand how the protocol operates, how it is expected to work, and identify any discrepancies or issues. Additionally, in the event of unexpected protocol changes, these logs provide a clear record of what has changed, making it easier to address breaking changes.

### Test Bench

The specification tests utilize predefined test benches to manage complex logic for starting, connecting, recording, and disconnecting agents in each test case. These test benches include:

- **Pair**: Lets to test the interaction between two nodes (node-node).
- **Line**: Allows to test a linear configuration of nodes (node-node-node).
- **Central Point**: Setups a central node connected to multiple nodes.
- **Loop**: Helps to test the behavior of nodes in a looped configuration.

By employing these test benches, the test code remains pretty clear and understandable, allowing developers to easily follow the logic of each test scenario while ensuring comprehensive coverage of different connection configurations.

## Development Environment

The development environment for DirectMQ is designed to provide a seamless experience for developers, utilizing a devcontainer defined by the standard devcontainer specifications with an Ubuntu operating system. This environment comes pre-installed with all necessary tools, requiring no additional configuration.

### Cloud and Local Usage

This development environment can be utilized with cloud services like GitHub Codespaces or locally with Devpod, allowing flexibility based on your working preferences.

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/sync-toys/DirectMQ)

### Supported Code Editor

The primary supported code editor is **Visual Studio Code** (VSCode) due to its popularity and capability to handle various programming languages and tools. VSCode is preconfigured to correctly locate all dependencies within the projects and includes defined debugging configurations for unit tests and specification agents.

### Devcontainer Features

The devcontainer also supports Devpod on Windows, cleverly circumventing performance issues with disk access in Windows Subsystem for Linux (WSL). This is achieved by mounting the workspace at `/host/workspace` and copying it to the native filesystem at `/workspace`, where you work within the devcontainer. 

However, there is a caveat: after a complete Docker reset, project files may disappear since they are no longer mounted in the host filesystem. Therefore, it's advisable to commit changes regularly to avoid data loss.

### Task Management

The development environment includes predefined scripts using `go-task` (taskfiles), streamlining various project operations so that developers do not need to memorize commands.

### Initial Setup

After launching the devcontainer, you can simply run the following command to build all necessary Docker images and compile all SDK project dependencies:

```bash
task prepare
```

### Running Tests

To run all unit and integration tests of SDKs, you can execute command:

```bash
task test
```

---

By following these guidelines, you can easily set up your development environment for DirectMQ, ensuring a productive and efficient workflow.

## License

DirectMQ is licensed under the MIT License. This means you are free to use, modify, and distribute the software, but it is provided "as is" without any warranties. 

### Summary of the MIT License

- **Permission**: You can use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the software.
- **Attribution**: You must include the original copyright notice and permission notice in all copies or substantial portions of the software.
- **No Warranty**: The software is provided "as is", without warranty of any kind, express or implied, including but not limited to the warranties of merchantability, fitness for a particular purpose, and non-infringement.

For more details, please refer to the full [MIT License](https://opensource.org/licenses/MIT).
