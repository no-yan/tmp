use oxc_allocator::Allocator;
use oxc_ast::AstKind;
use oxc_cfg::{
    EdgeType, Graph, Instruction, InstructionKind,
    graph::{
        Direction,
        graph::NodeIndex as CfgNodeIndex,
        visit::{Control, DfsEvent, EdgeRef, depth_first_search},
    },
};
use oxc_parser::Parser;
use oxc_semantic::{AstNode, AstNodes, SemanticBuilder};
use oxc_span::SourceType;

fn is_constructor(node: &AstNode) -> bool {
    matches!(node.kind(), AstKind::MethodDefinition(definition) if definition.kind.is_constructor())
}

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

    let nodes = semantic.semantic.nodes();
    let Some(class_node) = nodes.iter().find(|node| is_constructor(node)) else {
        return;
    };

    if !nodes
        .ancestor_kinds(class_node.id())
        .any(|kind| match kind {
            AstKind::Class(class) => class.super_class.is_some(),
            _ => false,
        })
    {
        println!("No super class extends");
        return;
    };

    let start = class_node.cfg_id();
    let cfg = semantic.semantic.cfg().unwrap();
    let _ = analyze(cfg, nodes, start);
}

const ERR_NO_SUPER_PATH: &str = "superが呼ばれずにconstructorが呼び出されるパスがあります!";

fn analyze(
    cfg: &oxc_cfg::ControlFlowGraph,
    nodes: &AstNodes,
    start: CfgNodeIndex,
) -> Result<(), String> {
    let graph = cfg.graph();
    // Classのcfgから、constructorを見つけます。
    let edge_reference = graph
        .edges_directed(start, Direction::Outgoing)
        .find(|edge_ref| matches!(edge_ref.weight(), EdgeType::NewFunction));

    let Some(edge_reference) = edge_reference else {
        return Ok(());
    };
    let maybe_constructor = edge_reference.target();

    let result = depth_first_search(graph, Some(maybe_constructor), |event| -> Control<String> {
        dbg!(event);
        match event {
            DfsEvent::Discover(basic_block_id, _) => {
                let super_instruction = cfg
                    .basic_block(basic_block_id)
                    .instructions()
                    .iter()
                    .find(|it| is_super_call_expression(it, nodes));

                let has_super = super_instruction.is_some();
                if has_super {
                    // super()が見つかったので、このパスはOK
                    // このパスの探索を打ち切る
                    Control::Prune
                } else if is_terminal_node(graph, basic_block_id) {
                    Control::Break(ERR_NO_SUPER_PATH.to_string())
                } else {
                    Control::Continue
                }
            }
            _ => Control::Continue,
        }
    });

    match result.break_value() {
        Some(val) => Err(val),
        None => Ok(()),
    }
}

fn is_super_call_expression(it: &Instruction, nodes: &AstNodes<'_>) -> bool {
    matches!(it.kind, InstructionKind::Statement if it.node_id.is_some_and(|node_id| {
        let node = nodes.get_node(node_id);
        matches!(&node.kind(), AstKind::ExpressionStatement(expr) if expr.expression.is_super_call_expression())
    }))
}

fn is_terminal_node(graph: &Graph, node_idx: CfgNodeIndex) -> bool {
    graph
        .neighbors_directed(node_idx, Direction::Outgoing)
        .count()
        == (0_usize)
}
#[cfg(test)]
mod tests {
    use oxc_allocator::Allocator;
    use oxc_parser::Parser;
    use oxc_semantic::SemanticBuilder;
    use oxc_span::SourceType;

    use crate::{ERR_NO_SUPER_PATH, analyze, is_constructor};

    #[test]
    fn test() {
        struct Case<'a> {
            code: &'a str,
            want: Result<(), String>,
        }

        let cases = vec![
            Case {
                code: r#"
class A extends B {
    constructor() {
        super()
    }
}"#,
                want: Ok(()),
            },
            Case {
                code: r#"
class A{
    constructor() {
        if (x) {
            return
        }
        super()
    }
}"#,
                want: Err(ERR_NO_SUPER_PATH.to_string()),
            },
            Case {
                code: r#"
class A{
    constructor() {
        if (x) {
            return
        } else {
            super()
        }
    }
}"#,
                want: Err(ERR_NO_SUPER_PATH.to_string()),
            },
        ];

        let source_ty = SourceType::default();
        for case in cases {
            let alloc = Allocator::default();
            let parsed = Parser::new(&alloc, case.code, source_ty).parse();
            assert_eq!(parsed.errors.is_empty(), true);

            let sema = SemanticBuilder::new().with_cfg(true).build(&parsed.program);
            assert_eq!(sema.errors.is_empty(), true);

            let cfg = sema.semantic.cfg().unwrap();
            let nodes = sema.semantic.nodes();
            let Some(class_node) = nodes.iter().find(|node| is_constructor(node)) else {
                unreachable!();
            };

            let ret = analyze(cfg, nodes, class_node.cfg_id());
            assert_eq!(ret, case.want);
        }
    }
}
