import generateRandomColorForWhiteText from "@/src/utils/generateRamdomColors"
import React, { useMemo } from "react"
import { Modal, Pressable, ScrollView, Text, View } from "react-native"

interface ColorsSelectorProps {
    visible: boolean
    onClose: () => void
    onSelect: (color: string) => void
}

const ColorsSelector: React.FC<ColorsSelectorProps> = ({ visible, onClose, onSelect }) => {
    const colors = useMemo(() => {
        return Array.from({ length: 100 }, () => generateRandomColorForWhiteText())
    }, [])
    return <Modal
        visible={visible}
        animationType="slide"
        presentationStyle="pageSheet"
        onRequestClose={onClose}
        style={{
            flex: 1
        }}
    >
        <ScrollView>

            <View className="flex-row flex-wrap gap-3 p-4 flex-1">
                {
                    colors.map((color, index) => (
                        <View className="w-1/3 min-w-[100px]" key={index}>
                            <Pressable
                                style={{ backgroundColor: color }}
                                onPress={() => { onSelect(color) }}
                                className="h-4 w-full py-4 px-2 rounded-lg"
                                key={index}
                            >
                                <Text className=" text-white">{color?.toUpperCase()}</Text>
                            </Pressable>
                        </View>
                    ))

                }
            </View>
        </ScrollView>
    </Modal>
}
export default ColorsSelector