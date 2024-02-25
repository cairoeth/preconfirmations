// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {ERC20} from "openzeppelin/token/ERC20/ERC20.sol";

contract ERC20Mock is ERC20 {
    constructor() ERC20("MockToken", "MOCK") {}

    function mint(address a, uint256 b) public {
        _mint(a, b);
    }
}
