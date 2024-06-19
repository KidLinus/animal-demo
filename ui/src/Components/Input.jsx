import { Flex, Text, Input as ChackraInput, useColorMode } from "@chakra-ui/react"

const Input = ({ value = "", onChange = () => { }, label, nullable, required, isDisabled, type, min, max, step, noBorderRight, ...props }) => {
    const { colorMode } = useColorMode()
    return <Flex pos="relative" w="full" {...props}>
        {label !== undefined && <Text zIndex="1" pos="absolute" top="0px" left="10px" fontSize="xs" lineHeight="0"
            borderBottom="1px solid" borderColor={colorMode == "light" ? "white" : "gray.700"} pl="1" pr="1">{label}</Text>}
        {required && <Text zIndex="1" pos="absolute" top="0px" right="10px" fontSize="xs" lineHeight="0"
            borderBottom="1px solid" borderColor={colorMode == "light" ? "white" : "gray.700"} pl="1" pr="1">required</Text>}
        <ChackraInput {...{ required, type, min, max, step, isDisabled }} value={value == null ? "" : value}
            {...(noBorderRight ? { borderRight: 0, borderRightRadius: 0 } : {})}
            onChange={e => onChange(nullable ? (e.target.value == "" ? null : e.target.value) : e.target.value)} />
    </Flex>
}

export default Input
