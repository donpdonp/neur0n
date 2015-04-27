class MasterControlProgram
  def initialize
    @machines = {}
  end

  def dispatch(msg)
    puts "admin.rb dispatch #{msg.inspect}"
    if msg['type'] == 'vm.add'
      if msg['name']
        machine = { id: newId, name: msg['name'] }
        puts "Adding machine #{machine}"
        add_machine(machine, msg)
      end
    end
    if msg['type'] == 'vm.list'
      #{machines: Neur0n::machine_list}
      puts "list #{@machines.inspect}"
      @machines
    end
  end

  def newId
    alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
    name = ""
    16.times{ name += alphabet[rand(36)]}
    name
  end

  def add_machine(machine, msg)
        idx = Neur0n::machine_add(machine[:id])
        if idx && msg["url"]
          @machines[machine[:id]] = machine
          url = machine[:url] = msg['url']
          puts "parsing #{url}"
          if gist_id = gistId(url)
            puts "gist id #{gist_id}"
            url = gist_api(gist_id)
          end
          puts "loading #{url}"
          code = Neur0n::http_get(url)
          Neur0n::machine_eval(machine[:id], code)
        end
  end

  def gistId(url)
    gist = url.match(/\/\/gist.github.com\/.*\/(.*)$/)
    return gist[1] if gist
  end

  def gist_api(id)
    gist_api = "https://api.github.com/gists/"+id
    gist = JSON.parse(Neur0n::http_get(gist_api))
    filename = gist['files'].keys.first
    return gist['files'][filename]['raw_url']
  end

end

MCP = MasterControlProgram.new

module Neur0n
  def self.dispatch(msg)
    MCP.dispatch(msg)
  end
end
