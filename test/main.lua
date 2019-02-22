local img, txt

function load()
  img = gfx.newimage("icon.png")
  txt = gfx.newtext(gfx.getfont(), "This is text")
end

function oninput(device, button, action, modifiers)
  if device == "keyboard" and button == "escape" and action == "release" then
    quit()
  end
end

function update(dt)
end

function draw()
  gfx.setcolor(1, 1, 1, 1)
  gfx.print({{"fps: ", getfps()}, {{1, 1, 0}, {0, 1, 1}}}, 0, 0)
  gfx.rectangle("fill", 300, 300, 480, 440)
  img:draw(300, 300)
  txt:draw(100, 100)
  gfx.setlinewidth(2)
  gfx.setlinejoin("bevel")
  gfx.ellipse("line", 300, 300, 100, 200)
  gfx.setcolor(1, 0, 0, 1)
  gfx.rectangle("line", 50, 50, 100, 100)
  gfx.setcolor(1, 1, 1, 1)
  gfx.line(0, 0, 100, 100, 200, 100)
end
