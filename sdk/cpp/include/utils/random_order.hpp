#pragma once

#include <list>
#include <random>

namespace directmq::utils {
template <typename T>
std::list<T> randomOrder(const std::list<T>& input) {
    std::list<T> result;
    result.insert(result.begin(), input.begin(), input.end());
    result.sort([](const T& a, const T& b) {
        (void)a;
        (void)b;
        return std::rand() % 2 == 0;
    });
    return result;
}
}  // namespace directmq::utils
