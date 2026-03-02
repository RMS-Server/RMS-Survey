import JSEncrypt from 'jsencrypt'

/**
 * Encrypt password using RSA public key
 * The backend returns raw base64 public key, we need to wrap it in PEM format
 */
export function encryptPassword(password: string, publicKeyBase64: string): string {
  const rsa = new JSEncrypt({})
  // Convert raw base64 to PEM format
  const pem = `-----BEGIN PUBLIC KEY-----\n${publicKeyBase64}\n-----END PUBLIC KEY-----`
  rsa.setPublicKey(pem)
  const encrypted = rsa.encrypt(password)
  return encrypted || password
}
