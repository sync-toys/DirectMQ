#pragma once

#include <list>
#include <memory>
#include <string>

#include "includes_string.hpp"

namespace directmq::utils {
std::list<std::shared_ptr<const std::string>> onlyUniqueStrings(
    const std::list<std::shared_ptr<const std::string>>& input) {
    std::list<std::shared_ptr<const std::string>> result;

    for (auto& str : input) {
        if (includesString(result, *str)) {
            continue;
        }

        result.push_back(str);
    }

    return result;
}
}  // namespace directmq::utils
