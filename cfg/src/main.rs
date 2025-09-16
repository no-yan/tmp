use oxc_allocator::Allocator;
use oxc_parser::Parser;
use oxc_semantic::SemanticBuilder;
use oxc_span::SourceType;

fn main() {
    let alloc = Allocator::default();
    let src_ty = SourceType::default();

    let code = std::fs::read_to_string("test/index.js").unwrap();
    let ret = Parser::new(&alloc, &code, src_ty).parse();
    if !ret.errors.is_empty() {
        for err in ret.errors {
            eprintln!("{err}");
        }
        return;
    }

    let program = ret.program;

    let semantic = SemanticBuilder::new()
        .with_check_syntax_error(true)
        .with_cfg(true)
        .build(&program);

    if !semantic.errors.is_empty() {
        for err in ret.errors {
            eprintln!("{err}");
        }
        return;
    }

    let cfg = semantic.semantic.cfg().unwrap();
    println!("{cfg:?}");
}
