import { Flex, Text, Select as ChakraSelect, useColorMode } from "@chakra-ui/react"

const Select = ({ value = "", onChange = () => { }, label, required, defaultValue, children, ...props }) => {
    const { colorMode } = useColorMode()
    return <Flex pos="relative" w="full" {...props}>
        {label !== undefined && <Text zIndex="1" pos="absolute" top="0px" left="10px" fontSize="xs" lineHeight="0"
            borderBottom="1px solid" borderColor={colorMode == "light" ? "white" : "gray.700"} pl="1" pr="1">{label}</Text>}
        {required && <Text zIndex="1" pos="absolute" top="0px" right="10px" fontSize="xs" lineHeight="0"
            borderBottom="1px solid" borderColor={colorMode == "light" ? "white" : "gray.700"} pl="1" pr="1">required</Text>}
        <ChakraSelect {...{ value, required, defaultValue }} onChange={e => onChange(e.target.value)}>{children}</ChakraSelect>
    </Flex>
}

export default Select
