use bump_scope::{BumpScope, BumpString, BumpVec};
use crate::list::ListVisitor;
use crate::r#struct::StructVisitor;
use crate::value::ArgValue;

pub trait Decoder<'a> {
    fn decode_i32(&mut self) -> Result<i32, DecodeError>;
    fn decode_borrowed_str(&mut self) -> Result<Result<&'a str, BumpString<'a, 'a>>, DecodeError>;
    fn decode_struct<V: StructVisitor<'a>>(&mut self, visitor: &mut V) -> Result<(), DecodeError>;
    fn decode_list<T, V: ListVisitor<'a, T>>(&mut self, visitor: &mut V) -> Result<(), DecodeError>;
    fn scope(&self) -> &'a BumpScope<'a>;

    #[cfg(feature = "std")]
    fn decode_owned_str(&mut self) -> Result<alloc::string::String, DecodeError>;
}

pub enum DecodeError {}

// pub trait DecodeHelper<'a>: Default {
//     type Value;
//     type MemoryHandle;
//
//     fn finish(self) -> (Self::Value, Option(Self::MemoryHandle));
// }
//
// impl<'a> DecodeHelper<'a> for i32 {
//     type Value = i32;
//     type MemoryHandle = ();
//
//     fn finish(self) -> (Self::Value, Option(Self::MemoryHandle)) {
//         (self, None)
//     }
// }
//
// #[derive(Default)]
// pub struct BorrowedStrHelper<'a> {
//     pub(crate) s: &'a str,
//     pub(crate) owner: Option<BumpString<'a, 'a>>,
// }
//
// impl<'a> DecodeHelper<'a> for BorrowedStrHelper<'a> {
//     type Value = &'a str;
//     type MemoryHandle = Option<BumpString<'a, 'a>>;
//
//     fn finish(self) -> (Self::Value, Self::MemoryHandle) {
//         (self.s, self.owner)
//     }
// }
//
//
// pub struct SliceHelper<'a, T: ArgValue<'a>> {
//     // TODO maybe there's a way that the underlying data could already be a slice so we can just borrow:
//     // pub(crate) s: &'a [T],
//     // TODO why not MutBumpVec?
//     pub(crate) vec: Option<BumpVec<'a, 'a, T>>,
//     pub(crate) helpers: BumpVec<'a, 'a, T::DecodeState::MemoryHandle>,
// }
//
// impl<'a, T: ArgValue<'a>> Default for SliceHelper<'a, T> {
//     fn default() -> Self {
//         Self {
//             vec: None,
//         }
//     }
// }
//
// impl<'a, T: ArgValue<'a>> DecodeHelper<'a> for SliceHelper<'a, T> {
//     type Value = &'a [T];
//     type MemoryHandle = (BumpVec<'a, 'a, T>, BumpVec<'a, 'a, T::DecodeState::MemoryHandle>);
//
//     fn finish(self) -> (Self::Value, Self::MemoryHandle) {
//         todo!()
//     }
// }
//
// impl<'a, T: ArgValue<'a>> ListVisitor<'a, T> for SliceHelper<'a, T> {
//     fn init(&mut self, len: usize, scope: &'a mut BumpScope<'a>) -> Result<(), DecodeError> {
//         let mut vec = BumpVec::new_in(scope);
//         vec.reserve(len);
//         Ok(())
//     }
//
//     fn next<D: Decoder<'a>>(&mut self, decoder: &mut D) -> Result<(), DecodeError> {
//         let vec = if let Some(vec) = &mut self.vec {
//             vec
//         } else {
//             let mut vec= BumpVec::new_in(decoder.scope());
//             self.vec = Some(vec);
//             self.vec.as_mut().unwrap()
//         };
//         let helper: T::DecodeState = Default::default();
//         helper.decode(decoder)?;
//         let (value, memory_handle) = helper.finish();
//         vec.push(value);
//     }
// }