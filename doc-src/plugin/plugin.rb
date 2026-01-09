require 'yard'
require 'yard-go'
require 'json'

module GoLinksHelper
  def signature(obj, link = true, show_extras = true, full_attr_name = true)
    case obj
    when YARDGo::CodeObjects::FuncObject
      if link && obj.has_tag?(:service_operation)
        ret = signature_types(obj, !link)
        args = obj.parameters.map { |m| m[0].split(/\s+/).last }.join(", ")
        line = "<strong>#{obj.name}</strong>(#{args}) #{ret}"
        return link ? linkify(obj, line) : line
      end
    end
    super(obj, link, show_extras, full_attr_name)
  end

  def html_syntax_highlight(source, type = nil)
    src = super(source, type || :go)
    object.has_tag?(:service_operation) ? link_types(src) : src
  end
end

YARD::Templates::Helpers::HtmlHelper.send(:prepend, GoLinksHelper)
YARD::Templates::Engine.register_template_path(File.dirname(__FILE__) + '/templates')

YARD::Parser::SourceParser.after_parse_list do
  YARD::Registry.all(:struct).each do |obj|
    if obj.file =~ /\/?service\/(.+?)\/(service|api)\.go$/
      obj.add_tag YARD::Tags::Tag.new(:service, $1)
      obj.groups = ["Constructor Functions", "Service Operations", "Request Methods", "Pagination Methods"]
    end
  end

  YARD::Registry.all(:method).each do |obj|
    begin
      if obj.file =~ /service\/.+?\/api\.go$/ && obj.scope == :instance
        if obj.name.to_s =~ /Pages$/
          obj.group = "Pagination Methods"
          opname = obj.name.to_s.sub(/Pages$/, '')
          obj.docstring = <<-eof
#{obj.name} iterates over pages of a {#{opname} #{opname}()} operation.
Return `false` in callback to stop iteration.
@example Iterating over at most 3 pages of a #{opname} operation
  pageNum := 0
  err := client.#{obj.name}(params, func(page *#{obj.parent.parent.name}.#{obj.parameters[1][0].split("*").last}, lastPage bool) bool {
    pageNum++
    fmt.Println(page)
    return pageNum <= 3
  })
@see #{opname}
eof
          obj.add_tag YARD::Tags::Tag.new(:paginator, '')
        elsif obj.name.to_s =~ /Request$/
          obj.group = "Request Methods"
          obj.signature = obj.name.to_s
          obj.parameters = []
          opname = obj.name.to_s.sub(/Request$/, '')
          obj.docstring = <<-eof
#{obj.name} generates a {aws/request.Request} object for {#{opname} #{opname}()}.
@example Sending a request using #{obj.name}() method
  req, resp := client.#{obj.name}(params)
  err := req.Send()
  if err == nil {
    fmt.Println(resp)
  }
eof
          obj.add_tag YARD::Tags::Tag.new(:request_method, '')
        else
          obj.group = "Service Operations"
          obj.add_tag YARD::Tags::Tag.new(:service_operation, '')
        end
      end
    rescue => e
      YARD::Logger.instance.warn "Skipping method #{obj.name} due to error: #{e}"
    end
  end

  begin
    apply_docs
  rescue => e
    YARD::Logger.instance.warn "Skipping apply_docs due to error: #{e}"
  end
end

# Safe version of apply_docs
def apply_docs
  svc_pkg = YARD::Registry.at('service')
  return if svc_pkg.nil?

  pkgs = svc_pkg.children.select { |t| t.type == :package }
  pkgs.each do |pkg|
    svc = pkg.children.find { |t| t.has_tag?(:service) }
    begin
      ctor = P(svc, ".New")
      svc_name = ctor.source[/ServiceName:\s*"(.+?)",/, 1]
      api_ver  = ctor.source[/APIVersion:\s*"(.+?)",/, 1]
    rescue YARD::CodeObjects::ProxyMethodError
      YARD::Logger.instance.warn "Skipping service #{svc&.name}, .New method not found"
      next
    end

    file = Dir.glob("models/apis/#{svc_name}/#{api_ver}/docs-2.json").sort.last
    next if file.nil?

    exmeth = svc.children.find { |s| s.has_tag?(:service_operation) }

    pkg.docstring += <<-eof
@example Using {#{svc.name}} client
  client := #{pkg.name}.New(nil)
  params := &#{pkg.name}.#{exmeth.parameters.first[0].split("*").last}{...}
  resp, err := client.#{exmeth.name}(params)
@see #{svc.name}
@version #{api_ver}
eof

    ctor.docstring += <<-eof
@example Constructing client with default config
  client := #{pkg.name}.New(nil)
@example Constructing client with custom config
  config := aws.NewConfig().WithRegion("us-west-2")
  client := #{pkg.name}.New(config)
eof

    json = JSON.parse(File.read(file))
    apply_doc(svc, json["service"]) if svc
    json["operations"].each do |op, doc|
      if doc && obj = svc.children.find { |t| t.name.to_s.downcase == op.downcase }
        apply_doc(obj, doc)
      end
    end
    json["shapes"].each do |shape, data|
      shape = shape_name(shape)
      if obj = pkg.children.find { |t| t.name.to_s.downcase == shape.downcase }
        apply_doc(obj, data["base"])
      end
      if data["refs"]
        data["refs"].each do |refname, doc|
          refshape, member = *refname.split("$")
          refshape = shape_name(refshape)
          if refobj = pkg.children.find { |t| t.name.to_s.downcase == refshape.downcase }
            if m = refobj.children.find { |t| t.name.to_s.downcase == member.downcase }
              apply_doc(m, doc || data["base"])
            end
          end
        end
      end
    end
  end
end

def apply_doc(obj, doc)
  tags = obj.docstring.tags || []
  obj.docstring = clean_docstring(doc)
  tags.each { |t| obj.docstring.add_tag(t) }
end

def shape_name(shape)
  shape.sub(/Request$/, "Input").sub(/Response$/, "Output")
end

def clean_docstring(docs)
  return nil unless docs
  docs = docs.gsub(/<!--.*?-->/m, '')
  docs = docs.gsub(/<fullname>.+?<\/fullname?>/m, '')
  docs = docs.gsub(/<examples?>.+?<\/examples?>/m, '')
  docs = docs.gsub(/<note>\s*<\/note>/m, '')
  docs = docs.gsub(/<a>(.+?)<\/a>/, '\1')
  docs = docs.gsub(/<note>(.+?)<\/note>/m) { "<div class=\"note\"><strong>Note:</strong> #{$1.gsub(/<\/?p>/, '')}</div>" }
  docs = docs.gsub(/\{(.+?)\}/, '`{\1}`')
  docs.gsub(/\s+/, ' ').strip
end
