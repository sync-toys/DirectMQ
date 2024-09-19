#pragma once

#include <list>
#include <memory>
#include <string>

namespace directmq::utils {
bool includesString(const std::list<std::shared_ptr<const std::string>>& input,
                    const std::string& str) {
    for (auto& item : input) {
        if (*item == str) {
            return true;
        }
    }

    return false;
}
}  // namespace directmq::utils
