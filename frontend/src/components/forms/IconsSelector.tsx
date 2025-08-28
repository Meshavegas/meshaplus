import generateRandomColorForWhiteText from "@/src/utils/generateRamdomColors"
import React, { useMemo } from "react"
import { Dimensions, Modal, Pressable, ScrollView, Text, View } from "react-native"
import { defaultIcons } from "@/src/utils/iconsData"
import Icon from "../Icons"

interface IconsSelectorProps {
    visible: boolean
    onClose: () => void
    onSelect: (color: string) => void
}

const IconsSelector: React.FC<IconsSelectorProps> = ({ visible, onClose, onSelect }) => {

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
                    defaultIcons.map((icon, index) => (
                        <View className="w-1/3" key={index}>
                            <Pressable
                                style={{ backgroundColor: generateRandomColorForWhiteText(),
                                    width: Dimensions.get('window').width / 3 - 20,
                                    height: Dimensions.get('window').width / 3 - 20,
                                    borderRadius: 10,
                                    justifyContent: 'center',
                                    alignItems: 'center',
                                 }}
                                onPress={() => { onSelect(icon.icon) }}
                                className="h-4 w-4 py-4 px-2 rounded-lg"

                                key={index}
                            >
                                <Icon name={icon.icon} size={40} color="white" />
                            </Pressable>
                        </View>
                    ))

                }
            </View>
        </ScrollView>
    </Modal>
}
export default IconsSelector