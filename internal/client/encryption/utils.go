package encryption

import (
	"encoding/json"
	"errors"

	userlib "github.com/cs161-staff/project2-starter-code/internal/userlib"
	"github.com/google/uuid"
)

func saveUser(storage StorageService, userdata *User) error {
	userBytes, err := json.Marshal(userdata)
	if err != nil {
		return err
	}

	encKey, err := userlib.HashKDF(userdata.MasterKey, []byte("enc"))
	if err != nil {
		return err
	}
	macKey, err := userlib.HashKDF(userdata.MasterKey, []byte("mac"))
	if err != nil {
		return err
	}

	payload, err := encryptAndMAC(userBytes, encKey, macKey)
	if err != nil {
		return err
	}
	hash := userlib.Hash([]byte(userdata.Username + "userStruct"))

	userUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		return err
	}

	if storage != nil {
		storage.Set(userUUID, payload)
	} else {
		userlib.DatastoreSet(userUUID, payload)
	}
	return nil
}

/* 加密数据并打包MAC封条的过程 */
func encryptAndMAC(data []byte, encKey []byte, macKey []byte) (payload []byte, err error) {
	aesKey := userlib.Hash(encKey)[:16]
	hmacKey := userlib.Hash(macKey)[:16]
	iv := userlib.RandomBytes(16)
	ciphertext := userlib.SymEnc(aesKey, iv, data)
	mac, err := userlib.HMACEval(hmacKey, ciphertext)
	if err != nil {
		return nil, err
	}

	payload = append(ciphertext, mac...)
	return payload, nil
}

func decryptAndVerify(payload []byte, encKey []byte, macKey []byte) (plaintext []byte, err error) {
	const macLen = 64

	if len(payload) < macLen {
		return nil, errors.New("malformed payload: data stream too short")
	}

	macOffset := len(payload) - macLen
	ciphertext := payload[:macOffset]
	receiveMac := payload[macOffset:]

	aesKey := userlib.Hash(encKey)[:16]
	hmacKey := userlib.Hash(macKey)[:16]

	expectMac, err := userlib.HMACEval(hmacKey, ciphertext)
	if err != nil {
		return nil, err
	}

	if !userlib.HMACEqual(receiveMac, expectMac) {
		return nil, errors.New("cryptographic doom: MAC verification failed, data tampered")
	}

	plaintext = userlib.SymDec(aesKey, ciphertext)

	return plaintext, nil
}

/* getPersonalKey: functions used encrypt Access Key Struct */
func getPersonalKey(masterKey []byte, filename string) (encKey []byte, macKey []byte) {
	salt := userlib.Hash([]byte(filename))
	baseKey, _ := userlib.HashKDF(masterKey, salt)

	encKey, _ = userlib.HashKDF(baseKey[:16], []byte("personal_enc"))
	macKey, _ = userlib.HashKDF(baseKey[:16], []byte("personal_mac"))
	return encKey, macKey
}

/* enxtends Globol File Key Pair */
func getFileKeys(fileKey []byte) (encKey []byte, macKey []byte) {
	encKey, _ = userlib.HashKDF(fileKey, []byte("file_enc"))
	macKey, _ = userlib.HashKDF(fileKey, []byte("file_mac"))
	return encKey[:16], macKey[:16]
}

/* extends Mailbox Key pair */
func getMailKeys(mailboxKey []byte) (mEncKey []byte, mMacKey []byte) {
	encKey, _ := userlib.HashKDF(mailboxKey, []byte("mailbox_enc"))
	macKey, _ := userlib.HashKDF(mailboxKey, []byte("mailbox_mac"))
	return encKey, macKey
}
