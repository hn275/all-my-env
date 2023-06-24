use hkdf;
use sha2::Sha256;
use std::env::var;

pub fn generate(components: &[String]) -> Result<[u8; 32], hkdf::InvalidLength> {
    let hkdf_secret = var("HKDF_SECRET").expect("`HKDF_SECRET` not set");
    let mut c: Vec<String> = Vec::from([hkdf_secret]);
    for i in components {
        c.push(i.to_owned());
    }

    let buf = components
        .into_iter()
        .map(|c| c.as_bytes())
        .collect::<Vec<&[u8]>>()
        .concat();

    let kdf = hkdf::Hkdf::<Sha256>::new(None, &buf);
    let mut output_key = [0u8; 32]; // hex encode: 16 -> 32 len
    kdf.expand(&[], &mut output_key)?;

    Ok(output_key)
}
